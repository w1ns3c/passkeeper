package tui

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"passkeeper/internal/config/client"
	"passkeeper/internal/logger"

	"passkeeper/internal/config"
	"passkeeper/internal/usecase/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	PageMain      = "Main"      // main form, here choose login/register action
	PageLogin     = "Login"     // form for user login
	PageRegister  = "Register"  // form for user registration
	PageCredsMenu = "CredsMenu" // form for user registration

	PageAuthError       = "AuthErr"      // auth error modal
	PageRegisterError   = "RegErr"       // register error modal
	PageCredsListErr    = "CredsListErr" //creds list error modal
	PageCredUpdError    = "CredUpdErr"
	PageRegisterSuccess = "RegSuccess" // register success modal
	PageAuthed          = "Authed"     //login success modal

	SubPageBank  = "Banking"
	SubPageNotes = "Notes"
	SubPageCreds = "Resources" // full form for view/edit/add user cred

	Hint = tview.NewTextView().
		SetTextColor(tcell.ColorBisque).
		SetText("(Esc) to back to main page")

	HintTextCreds = "" +
		"(Space/Enter) - to choose credential\n" +
		"(a)           - to add new credential\n" +
		"(e)           - to edit credential\n" +
		"(d/Del)       - to delete credential\n" +
		"(Ctrl+L)      - to logout\n"
	HintCreds = tview.NewTextView().
			SetTextColor(tcell.ColorBisque).
			SetText(HintTextCreds)

	PassHidden = "******"
)

// TUI struct is tui client app
type TUI struct {
	App        *tview.Application
	Screen     tcell.Screen
	Pages      *tview.Pages // main menu | login/register user
	SubPages   *tview.Pages // cred blobs pages
	MinPassLen int

	// Actions
	Usecase   cli.ClientUsecase
	Ctx       context.Context
	ctxCancel context.CancelFunc
	wg        *sync.WaitGroup

	//Creds []*entities.Credential

	// logger
	log *zerolog.Logger

	// save TUI forms
	FormAuth      *tview.Flex
	FormLogin     *tview.Flex
	FormReg       *tview.Flex
	FormCredsMenu *tview.Flex

	// save TUI SubPage forms
	SubformCreds *tview.Flex
	SubformBank  *tview.Flex
	SubformNotes *tview.Flex
}

// NewTUIconf is wrapper for NewTUI constructor with config parser
func NewTUIconf(conf *client.Args) (tui *TUI, err error) {
	return NewTUI(conf.Addr, conf.LogLevel, conf.LogFile, conf.SyncTime)
}

// NewTUI func is constructor for TUI
func NewTUI(addr string, debugLevel, logFile string, syncTime int) (tui *TUI, err error) {
	lg := logger.InitFile(debugLevel, logFile)

	scr, err := tcell.NewScreen()
	if err != nil {
		lg.Error().Err(err).
			Msg("failed to create tcell screen")

		return nil, err
	}
	pages := tview.NewPages()
	subPages := tview.NewPages()

	usecase, err := cli.NewClientUC(
		cli.WithAddr(addr),
		cli.WithSyncTime(time.Duration(syncTime)*time.Second))
	if err != nil {
		lg.Error().Err(err).
			Msg("failed to create client usecase")

		return nil, err
	}

	tuiApp := &TUI{

		Screen:   scr,
		Pages:    pages,
		SubPages: subPages,
		App: tview.NewApplication().
			SetScreen(scr).SetRoot(pages, true).
			EnableMouse(false),

		MinPassLen: config.MinPassLen,
		Usecase:    usecase,
		wg:         &sync.WaitGroup{},

		log: lg,
	}

	// init pages
	tuiApp.FormAuth = FlexAuth(tuiApp)
	tuiApp.FormLogin = NewLoginForm(tuiApp)
	tuiApp.FormReg = NewRegisterForm(tuiApp)
	authedForm := NewModalWithParams(tuiApp, "Success!", PageCredsMenu)

	header := NewHeader(0)
	tuiApp.FormCredsMenu = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(subPages, 0, 15, true)

	// Main pages
	// add pages
	tuiApp.Pages.AddPage(PageMain, tuiApp.FormAuth, true, true) // login/register select form
	tuiApp.Pages.AddPage(PageLogin, tuiApp.FormLogin, true, false)
	tuiApp.Pages.AddPage(PageRegister, tuiApp.FormReg, true, false)
	tuiApp.Pages.AddPage(PageAuthed, authedForm, true, false)
	tuiApp.Pages.AddPage(PageCredsMenu, tuiApp.FormCredsMenu, true, false)
	//tuiApp.Pages.AddPage(SubPageCreds, credsForm, true, false)
	//tuiApp.FormMain = NewMainMenu(tuiApp.Pages)

	// Subpages
	//tuiApp.SubPages.AddPage(SubPageCreds, subPages, true, false)
	//tuiApp.SubPages.AddPage(SubPageBank, subPages, true, false)
	//tuiApp.SubPages.AddPage(SubPageNotes, subPages, true, false)

	// change hints layout
	HintCreds.SetBorderPadding(0, 0, 1, 0)
	Hint.SetBorderPadding(0, 0, 1, 0)

	return tuiApp, nil
}

func (tuiApp *TUI) Logout() error {
	tuiApp.Ctx = context.Background()
	tuiApp.Usecase.Logout()

	return nil
}

func (tuiApp *TUI) Run(ctx context.Context) error {
	tuiApp.Ctx = ctx

	tuiApp.log.Info().
		Msg("starting client")

	tuiApp.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			err := tuiApp.interruptSignal()
			if err != nil {
				tuiApp.log.Error().
					Err(err).Msg("cannot send interrupt signal")

				return event
			}

		case tcell.KeyF2:
			if tuiApp.Usecase.IsAuthed() {
				tuiApp.SubformCreds = NewCredsList(tuiApp)
				item := tuiApp.FormCredsMenu.GetItem(0).(*Header)
				item.ChangePage(1)
				tuiApp.SubPages.AddPage(SubPageCreds, tuiApp.SubformCreds, true, false)
				tuiApp.SubPages.SwitchToPage(SubPageCreds)
			}

		case tcell.KeyF3:
			if tuiApp.Usecase.IsAuthed() {
				tuiApp.SubformBank = NewBanking(tuiApp.Usecase.GetCards())
				item := tuiApp.FormCredsMenu.GetItem(0).(*Header)
				item.ChangePage(2)
				tuiApp.SubPages.AddPage(SubPageBank, tuiApp.SubformBank, true, false)
				tuiApp.SubPages.SwitchToPage(SubPageBank)
			}
		case tcell.KeyF4:
			if tuiApp.Usecase.IsAuthed() {
				tuiApp.SubformNotes = NewNotes(tuiApp.Usecase.GetNotes())
				item := tuiApp.FormCredsMenu.GetItem(0).(*Header)
				item.ChangePage(3)
				tuiApp.SubPages.AddPage(SubPageNotes, tuiApp.SubformNotes, true, false)
				tuiApp.SubPages.SwitchToPage(SubPageNotes)
			}
		default:
			return event
		}

		return event
	})

	go func() {
		err := tuiApp.App.Run()
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to start app")
			os.Exit(1)
		}
	}()

	// Block the rest of the code until a signal is received.
	sig := <-ctx.Done()
	tuiApp.log.Info().Str("sig", fmt.Sprintf("%v", sig)).Msg("Got signal")
	tuiApp.log.Info().Msg("Shutting everything down gracefully")

	tuiApp.Stop()

	return nil

	//return err
}

func (tuiApp *TUI) Stop() {
	// wait async creds update
	tuiApp.wg.Wait()

	tuiApp.log.Info().Msg("[i] Client stopped correctly!")

}

func (tuiApp *TUI) interruptSignal() error {
	tuiApp.log.Info().
		Msg("catch \"Ctrl+C\" signal")

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		tuiApp.log.Error().Err(err).
			Msg("failed to find app process, can't notify context")

		return err
	}

	if err := p.Signal(os.Interrupt); err != nil {
		tuiApp.log.Error().Err(err).
			Msg("failed to send signal to app process, can't notify context")

		return err
	}

	return nil
}

func genHelp(page string) string {
	prefix := " "
	return fmt.Sprintf(""+
		prefix+"(Space/Enter) - to choose %s\n"+
		prefix+"(a)           - to add new %s\n"+
		prefix+"(e)           - to edit %s\n"+
		prefix+"(d/Del)       - to delete %s\n"+
		prefix+"(Ctrl+L)      - to logout\n", page, page, page, page)

}
