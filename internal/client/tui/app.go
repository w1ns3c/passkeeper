package tui

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"passkeeper/internal/config/client"
	"passkeeper/internal/logger"

	"passkeeper/internal/config"
	"passkeeper/internal/usecase/cli"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	PageMain            = "Auth"
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
	Usecase cli.ClientUsecase
	Ctx     context.Context

	//Creds []*entities.Credential

	// logger
	log *zerolog.Logger
}

// NewTUIconf is wrapper for NewTUI constructor with config parser
func NewTUIconf(conf *client.Args) (tui *TUI, err error) {
	return NewTUI(conf.Addr, conf.LogLevel)
}

// NewTUI func is constructor for TUI
func NewTUI(addr string, debugLevel string) (tui *TUI, err error) {
	scr, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	pages := tview.NewPages()

	//ctx := context.Background()
	usecase, err := cli.NewClientUC(addr)
	if err != nil {
		return nil, err
	}

	tuiApp := &TUI{

		Screen: scr,
		Pages:  pages,
		App: tview.NewApplication().
			SetScreen(scr).SetRoot(pages, true).
			EnableMouse(false),
		//Creds:      make([]*entities.Credential, 0),
		MinPassLen: config.MinPassLen,
		Usecase:    usecase,

		Ctx: context.Background(),
		log: logger.Init(debugLevel),
	}

	//tuiApp.Creds = CredsList

	var (
		// init pages
		mainForm  = FlexMain(tuiApp)
		loginForm = NewLoginForm(tuiApp)
		regForm   = NewRegisterForm(tuiApp)
		//credsForm  = NewCredsList(tuiApp)
		authedForm = NewModalWithParams(tuiApp, "Success!", PageCreds)
	)

	// add pages
	tuiApp.Pages.AddPage(PageMain, mainForm, true, true)
	tuiApp.Pages.AddPage(PageLogin, loginForm, true, false)
	tuiApp.Pages.AddPage(PageRegister, regForm, true, false)
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
