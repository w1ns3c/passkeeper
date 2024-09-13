package tui

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"passkeeper/internal/config/client"
	"passkeeper/internal/logger"
	"sync"
	"syscall"
	"time"

	"passkeeper/internal/config"
	"passkeeper/internal/usecase/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	PageMain            = "Main"
	PageLogin           = "Login"
	PageRegister        = "Register"
	PageCreds           = "Credentials"
	PageAuthError       = "AuthErr"
	PageRegisterError   = "RegErr"
	PageCredsListErr    = "CredsListErr"
	PageCredUpdError    = "CredUpdErr"
	PageRegisterSuccess = "RegSuccess"
	PageAuthed          = "Authed"

	TmpPage = "tmppage"

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

	ErrNotAuthed = fmt.Errorf("not authed")
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

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	tuiApp := &TUI{

		Screen: scr,
		Pages:  pages,
		App: tview.NewApplication().
			SetScreen(scr).SetRoot(pages, true).
			EnableMouse(false),
		//Creds:      make([]*entities.Credential, 0),
		MinPassLen: config.MinPassLen,
		Usecase:    usecase,
		wg:         &sync.WaitGroup{},

		Ctx:       ctx,
		ctxCancel: stop,

		log: lg,
	}

	//tuiApp.Creds = CredsList

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
	//tui.Ctx = context.Background()
	tui.Usecase.Logout()

	return nil
}

func (tui *TUI) Run() error {
	tui.log.Info().
		Msg("starting client")

	err := tui.App.Run()
	if err != nil {
		tui.log.Error().
			Err(err).Msg("failed to start app")
	}

	return err
}

func (tui *TUI) Stop() {
	tui.ctxCancel()

	// wait async creds update
	tui.wg.Wait()

	tui.log.Info().Msg("[i] Client stopped correctly!")

}
