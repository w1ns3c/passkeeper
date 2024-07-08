package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/w1ns3c/passkeeper/internal/entities"
	"net/mail"
	"regexp"
	"strings"
)

const MinPassLen = 8

var (
	PageMain            = "Auth"
	PageLogin           = "Login"
	PageRegister        = "Register"
	PageCreds           = "Credentials"
	PageAuthErr         = "AuthErr"
	PageRegisterError   = "RegErr"
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
)

type TUI struct {
	App        *tview.Application
	Screen     tcell.Screen
	Pages      *tview.Pages
	MinPassLen int

	// user info
	User  *entities.User
	Creds []entities.Credential
}

func NewTUI() (tui *TUI, err error) {
	scr, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	pages := tview.NewPages()

	tuiApp := &TUI{
		Screen: scr,
		Pages:  pages,
		App: tview.NewApplication().
			SetScreen(scr).SetRoot(pages, true).
			EnableMouse(false),
		Creds:      make([]entities.Credential, 0),
		MinPassLen: MinPassLen,
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

func (app *TUI) GetCreds() (creds []entities.Credential, err error) {
	if app.User == nil {
		return nil, fmt.Errorf("not authed")
	}

	if app.User.Login == "user" {
		return CredsList, nil
	}

	if app.User.Login == "error" {
		return nil, fmt.Errorf("can't get credentials")
	}

	return make([]entities.Credential, 0), nil
}

func (app *TUI) Logout() error {
	app.Creds = nil
	app.User = nil

	return nil
}

func (app *TUI) Register(email, username, password, passRepeat string) error {
	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)

	err := app.FilterUserRegValues(email, username, password, passRepeat)
	if err != nil {
		return err
	}

	// TODO Change this
	users[username] = password

	return nil
}

func (app *TUI) FilterUserRegValues(email, username, password, passRepeat string) error {
	if username == "" {
		return fmt.Errorf("username is empty")
	}
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email is not valid")
	}
	// from here https://emaillistvalidation.com/blog/mastering-email-validation-in-golang-crafting-robust-regex-patterns/
	if result, _ := regexp.MatchString("^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", email); !result {
		return fmt.Errorf("email is not valid by new regexp")
	}

	if password != passRepeat {
		return fmt.Errorf("passwords are not the same")
	}

	// TODO Uncomment it
	//if len(password) < app.MinPassLen {
	//	return fmt.Errorf("password len should be a least %d signs", app.MinPassLen)
	//}

	if _, ok := users[username]; ok {
		return fmt.Errorf("user already exist")
	}

	return nil
}

func (app *TUI) Login(username, password string) bool {
	pass, ok := users[username]
	if !ok {
		return false
	}

	return pass == password
}
