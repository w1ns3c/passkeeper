package tui

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"passkeeper/internal/config/client"
	"passkeeper/internal/logger"
	"sync"
	"time"

	"passkeeper/internal/config"
	"passkeeper/internal/usecase/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	PageMain            = "Main"         // main form, here choose login/register action
	PageLogin           = "Login"        // form for user login
	PageRegister        = "Register"     // form for user registration
	PageCreds           = "Credentials"  // full form for view/edit/add user cred
	PageAuthError       = "AuthErr"      // auth error modal
	PageRegisterError   = "RegErr"       // register error modal
	PageCredsListErr    = "CredsListErr" //creds list error modal
	PageCredUpdError    = "CredUpdErr"
	PageRegisterSuccess = "RegSuccess" // register success modal
	PageAuthed          = "Authed"     //login success modas

	Hint = tview.NewTextView().
		SetTextColor(tcell.ColorBisque).
		SetText("(Esc) to back to main page")

	HintText = "" +
		"(Space/Enter) - to choose credential\n" +
		"(a)           - to add new credential\n" +
		"(e)           - to edit credential\n" +
		"(d/Del)       - to delete credential\n" +
		"(Ctrl+L)      - to logout\n"
	HintCreds = tview.NewTextView().
			SetTextColor(tcell.ColorBisque).
			SetText(HintText)

	PassHidden = "******"
)

// TUI struct is tui client app
type TUI struct {
	App        *tview.Application
	Screen     tcell.Screen
	Pages      *tview.Pages
	MinPassLen int

	// Actions
	Usecase   cli.ClientUsecase
	Ctx       context.Context
	ctxCancel context.CancelFunc
	wg        *sync.WaitGroup

	//Creds []*entities.Credential

	// logger
	log *zerolog.Logger

	// save TUI primitives
	FormMain  *tview.Flex
	FormLogin *tview.Flex
	FormReg   *tview.Flex
	FormCreds *tview.Flex
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

	//ctx := context.Background()
	usecase, err := cli.NewClientUC(
		cli.WithAddr(addr),
		cli.WithSyncTime(time.Duration(syncTime)*time.Second))
	if err != nil {
		lg.Error().Err(err).
			Msg("failed to create client usecase")

		return nil, err
	}

	tuiApp := &TUI{

		Screen: scr,
		Pages:  pages,
		App: tview.NewApplication().
			SetScreen(scr).SetRoot(pages, true).
			EnableMouse(false),

		MinPassLen: config.MinPassLen,
		Usecase:    usecase,
		wg:         &sync.WaitGroup{},

		log: lg,
	}

	// init pages
	tuiApp.FormMain = FlexMain(tuiApp)
	tuiApp.FormLogin = NewLoginForm(tuiApp)
	tuiApp.FormReg = NewRegisterForm(tuiApp)
	//credsForm  = NewCredsList(tuiApp)
	authedForm := NewModalWithParams(tuiApp, "Success!", PageCreds)

	// add pages
	tuiApp.Pages.AddPage(PageMain, tuiApp.FormMain, true, true) // login/register select form
	tuiApp.Pages.AddPage(PageLogin, tuiApp.FormLogin, true, false)
	tuiApp.Pages.AddPage(PageRegister, tuiApp.FormReg, true, false)
	//tuiApp.Pages.AddPage(PageCreds, credsForm, true, false)
	tuiApp.Pages.AddPage(PageAuthed, authedForm, true, false)

	// change hints layout
	HintCreds.SetBorderPadding(0, 0, 1, 0)
	Hint.SetBorderPadding(0, 0, 1, 0)

	return tuiApp, nil
}

func (tui *TUI) Logout() error {
	tui.Ctx = context.Background()
	tui.Usecase.Logout()

	return nil
}

func (tui *TUI) Run(ctx context.Context) error {
	tui.Ctx = ctx

	tui.log.Info().
		Msg("starting client")

	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			err := tui.interruptSignal()
			if err != nil {
				tui.log.Error().
					Err(err).Msg("cannot send interrupt signal")

				return event
			}

		default:
			return event
		}

		return event
	})

	go func() {
		err := tui.App.Run()
		if err != nil {
			tui.log.Error().
				Err(err).Msg("failed to start app")
			os.Exit(1)
		}
	}()

	// Block the rest of the code until a signal is received.
	sig := <-ctx.Done()
	tui.log.Info().Str("sig", fmt.Sprintf("%v", sig)).Msg("Got signal")
	tui.log.Info().Msg("Shutting everything down gracefully")

	tui.Stop()

	return nil

	//return err
}

func (tui *TUI) Stop() {
	// wait async creds update
	tui.wg.Wait()

	tui.log.Info().Msg("[i] Client stopped correctly!")

}

func (tui *TUI) interruptSignal() error {
	tui.log.Info().
		Msg("catch \"Ctrl+C\" signal")

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		tui.log.Error().Err(err).
			Msg("failed to find app process, can't notify context")

		return err
	}

	if err := p.Signal(os.Interrupt); err != nil {
		tui.log.Error().Err(err).
			Msg("failed to send signal to app process, can't notify context")

		return err
	}

	return nil
}
