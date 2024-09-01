package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"passkeeper/internal/entities"
	"strings"
)

func FlexMain(tuiApp *TUI) *tview.Flex {

	btn1 := tview.NewButton("Login").
		SetSelectedFunc(func() {
			tuiApp.Pages.SwitchToPage(PageLogin)
		})

	btn1Flex := tview.NewFlex().
		AddItem(nil, 0, 2, false).
		AddItem(btn1, 5, 1, true).
		AddItem(nil, 0, 2, false).
		SetDirection(tview.FlexRow)

	btn2 := tview.NewButton("Register").
		SetSelectedFunc(func() {
			tuiApp.Pages.SwitchToPage(PageRegister)
		})

	btn2Flex := tview.NewFlex().
		AddItem(nil, 0, 2, false).
		AddItem(btn2, 5, 1, true).
		AddItem(nil, 0, 2, false).
		SetDirection(tview.FlexRow)

	flex := tview.NewFlex().
		AddItem(nil, 0, 2, false).
		AddItem(btn1Flex, 10, 1, true).
		AddItem(nil, 0, 1, false).
		AddItem(btn2Flex, 10, 1, true).
		AddItem(nil, 0, 2, false)
	flex.SetTitle("Pass Keeper").SetBorder(true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		item1 := flex.GetItem(1) // login button
		item2 := flex.GetItem(3) // register button
		switch event.Key() {
		default:
			//tuiApp.Pages.SwitchToPage(PageCreds)
			return event
		case tcell.KeyEsc:
			page := "exit"
			exitModal := ExitModal(tuiApp)
			tuiApp.Pages.AddPage(page, exitModal, true, false)
			tuiApp.Pages.SwitchToPage(page)

		case tcell.KeyTab:

			if !item1.HasFocus() &&
				!btn1Flex.HasFocus() && !btn1.HasFocus() {
				tuiApp.App.SetFocus(btn1)
				return event
			}

			if !item2.HasFocus() &&
				!btn2Flex.HasFocus() && !btn2.HasFocus() {
				tuiApp.App.SetFocus(btn2)
				return event
			}

		case tcell.KeyLeft:

			if !item1.HasFocus() &&
				!btn1Flex.HasFocus() && !btn1.HasFocus() {
				tuiApp.App.SetFocus(btn1)
				return event
			}

			if !item2.HasFocus() &&
				!btn2Flex.HasFocus() && !btn2.HasFocus() {
				tuiApp.App.SetFocus(btn2)
				return event
			}

		case tcell.KeyRight:

			if !item1.HasFocus() &&
				!btn1Flex.HasFocus() && !btn1.HasFocus() {
				tuiApp.App.SetFocus(btn1)
				return event
			}

			if !item2.HasFocus() &&
				!btn2Flex.HasFocus() && !btn2.HasFocus() {
				tuiApp.App.SetFocus(btn2)
				return event
			}

		case tcell.KeyEnter:
			if btn1.HasFocus() {
				tuiApp.Pages.SwitchToPage(PageLogin)
			}
			if btn2.HasFocus() {
				tuiApp.Pages.SwitchToPage(PageRegister)
			}
			return event

			// DEBUG
			//case tcell.KeyUp:
			//	tuiApp.Pages.SwitchToPage(PageCreds)
			//case tcell.KeyDown:
			//	details := NewDetailsForm(tuiApp)
			//	tuiApp.Pages.AddPage(TmpPage, details, true, false)
			//	tuiApp.Pages.SwitchToPage(TmpPage)

		}
		return event
	})

	return flex
}

type LoginStruct struct {
	Login string
	Pass  string
}

type RegisterStruct struct {
	Email  string
	Login  string
	Pass   string
	Repeat string
}

func NewLoginForm(tuiApp *TUI) *tview.Flex {

	var loginForm = tview.NewForm().
		SetItemPadding(1)

	user := LoginStruct{}

	loginForm.AddInputField("Username", "", 20, nil, func(login string) {
		user.Login = login
	})
	loginForm.AddInputField("Password", "", 20, nil, func(password string) {
		user.Pass = password
	})

	clearForm := func() {
		if loginForm.GetFormItemCount() < 2 {
			return
		}
		fieldLogin := loginForm.GetFormItem(0).(*tview.InputField)
		fieldPassword := loginForm.GetFormItem(1).(*tview.InputField)

		fieldLogin.SetText("")
		fieldPassword.SetText("")
		loginForm.SetFocus(0)
	}

	loginForm.AddButton("Login", func() {
		itemUser := loginForm.GetFormItem(0)
		uField := itemUser.(*tview.InputField)
		username := uField.GetText()

		itemPass := loginForm.GetFormItem(1)
		pField := itemPass.(*tview.InputField)
		password := pField.GetText()

		//authed := Login(username, password)
		err := tuiApp.Usecase.Login(tuiApp.ctx, username, password)
		if err != nil {
			// not authed
			errAuthForm := NewModalWithParams(tuiApp, "Wrong username/password!", PageLogin)
			tuiApp.Pages.AddPage(PageAuthErr, errAuthForm, true, false)
			tuiApp.Pages.SwitchToPage(PageAuthErr)
			return
		}

		// TODO change this auth
		tuiApp.User = &entities.User{
			ID:     "",
			Login:  username,
			Hash:   password,
			Secret: "",
			Phone:  "",
		}

		creds, err := tuiApp.GetCreds()
		if err != nil {
			PageCredsErr := "credsError"
			errModal := NewModalWithParams(tuiApp, err.Error(), PageLogin)
			tuiApp.Pages.AddPage(PageCredsErr, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredsErr)
			return
		}

		clearForm()

		tuiApp.Creds = creds
		credsForm := NewCredsList(tuiApp)
		tuiApp.Pages.AddPage(PageCreds, credsForm, true, false)

		tuiApp.Pages.SwitchToPage(PageAuthed)

	})

	loginForm.AddButton("Cancel", func() {
		clearForm()
		tuiApp.Pages.SwitchToPage(PageMain)
	})

	loginForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			fallthrough
		case tcell.KeyEsc:
			clearForm()
			tuiApp.Pages.SwitchToPage(PageMain)
		}

		return event
	})

	loginFlex := tview.NewFlex()
	loginFlex.SetDirection(tview.FlexRow).
		AddItem(loginForm, 0, 4, true).
		AddItem(Hint, 0, 1, false)

	return loginFlex
}

func NewRegisterForm(tuiApp *TUI) *tview.Flex {
	var regForm = tview.NewForm().
		SetItemPadding(1)

	regForm.AddInputField("Email", "", 20, nil, nil)
	regForm.AddInputField("Username", "", 20, nil, nil)
	regForm.AddInputField("Password", "", 20, nil, nil)
	regForm.AddInputField("Repeat pass", "", 20, nil, nil)

	clearForm := func() {
		if regForm.GetFormItemCount() < 4 {
			return
		}

		fieldEmail := regForm.GetFormItem(0).(*tview.InputField)
		fieldLogin := regForm.GetFormItem(1).(*tview.InputField)
		fieldPassword := regForm.GetFormItem(2).(*tview.InputField)
		fieldRepeat := regForm.GetFormItem(3).(*tview.InputField)

		fieldEmail.SetText("")
		fieldLogin.SetText("")
		fieldPassword.SetText("")
		fieldRepeat.SetText("")
		regForm.SetFocus(0)
	}

	regForm.AddButton("Register", func() {
		itemEmail := regForm.GetFormItem(0)
		emailField := itemEmail.(*tview.InputField)
		email := emailField.GetText()

		itemUser := regForm.GetFormItem(1)
		uField := itemUser.(*tview.InputField)
		username := uField.GetText()

		itemPass := regForm.GetFormItem(2)
		pField := itemPass.(*tview.InputField)
		password := pField.GetText()

		itemPass2 := regForm.GetFormItem(3)
		pField2 := itemPass2.(*tview.InputField)
		repeat := pField2.GetText()

		if err := tuiApp.Usecase.Register(tuiApp.ctx, email, username, password, repeat); err != nil {
			errModal := NewModalWithParams(tuiApp, err.Error(), PageRegister)
			tuiApp.Pages.AddPage(PageRegisterError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageRegisterError)

			return
		}
		clearForm()

		text := "User \"" + username + "\" registered!"
		registeredModal := NewModalWithParams(tuiApp, text, PageMain)
		tuiApp.Pages.AddPage(PageRegisterSuccess, registeredModal, true, false)
		tuiApp.Pages.SwitchToPage(PageRegisterSuccess)
	})
	regForm.AddButton("Cancel", func() {
		tuiApp.Pages.SwitchToPage(PageMain)
	})

	regForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			tuiApp.Pages.SwitchToPage(PageMain)
		}

		return event
	})

	registerFlex := tview.NewFlex()
	registerFlex.SetDirection(tview.FlexRow).
		AddItem(regForm, 0, 4, true).
		AddItem(Hint, 0, 1, false)

	return registerFlex
}

// FilterResource beautify resource name to print it in list
func FilterResource(res string) string {

	// filter protocol ( https:// & etc)
	parts := strings.Split(res, ":/")
	if len(parts) != 1 {
		res = parts[1]
	}

	// trim left "/" sign, if it exists
	res = strings.TrimLeft(res, "/")

	// filter uri site/<path_to_remove>
	parts = strings.Split(res, "/")
	res = parts[0]

	return res
}

// FilterEmail validate email address
func FilterEmail(email string) bool {
	// TODO rewrite this function
	return true
}
