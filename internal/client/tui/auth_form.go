package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"passkeeper/internal/config"
)

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

	clearForm := func(form *tview.Form) *tview.Form {
		form.Clear(false)
		form.AddInputField("Username", "", 20, nil, func(login string) {
			user.Login = login
		})
		form.AddInputField("Password", "", 20, nil, func(password string) {
			user.Pass = password
		})

		return form
	}

	loginForm.AddButton("Login", func() {
		itemUser := loginForm.GetFormItem(0)
		uField := itemUser.(*tview.InputField)
		username := uField.GetText()

		itemPass := loginForm.GetFormItem(1)
		pField := itemPass.(*tview.InputField)
		password := pField.GetText()

		err := tuiApp.Usecase.Login(tuiApp.Ctx, username, password)
		if err != nil {
			stErr := status.FromContextError(err)
			if strings.Contains(stErr.Message(), "connection refused") {
				tuiApp.log.Error().Err(err).Msg("server is unavailable")
				errAuthForm := NewModalWithParams(tuiApp, "Server is unavailable!", PageLogin)
				tuiApp.Pages.AddPage(PageAuthError, errAuthForm, true, false)
				tuiApp.Pages.SwitchToPage(PageAuthError)
				return
			}

			// TODO delete log output user/pass to log
			// not authed
			tuiApp.log.Error().Err(err).Msgf("wrong username or password: %s:%s", username, password)
			errAuthForm := NewModalWithParams(tuiApp, "Wrong username/password!", PageLogin)
			tuiApp.Pages.AddPage(PageAuthError, errAuthForm, true, false)
			tuiApp.Pages.SwitchToPage(PageAuthError)
			return
		}

		tuiApp.log.Error().Err(err).Msgf("User %s successfully logged in", username)
		md := metadata.New(map[string]string{config.TokenHeader: tuiApp.Usecase.GetToken()})
		tuiApp.Ctx = metadata.NewOutgoingContext(tuiApp.Ctx, md)

		err = tuiApp.Usecase.GetBlobs(tuiApp.Ctx)
		if err != nil {
			tuiApp.log.Error().Err(err).Msg("failed to get creds from server")
			errModal := NewModalWithParams(tuiApp, err.Error(), PageLogin)
			tuiApp.Pages.AddPage(PageCredsListErr, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredsListErr)

			return
		}

		f := tuiApp.FormLogin.GetItem(0).(*tview.Form)
		f = clearForm(f)

		tuiApp.SubformCreds = NewCredsList(tuiApp)
		tuiApp.Pages.SwitchToPage(PageAuthed)

		item := tuiApp.FormCredsMenu.GetItem(0).(*Header)
		item.ChangePage(1)
		tuiApp.SubPages.AddPage(SubPageCreds, tuiApp.SubformCreds, true, false)
		tuiApp.SubPages.SwitchToPage(SubPageCreds)

		tuiApp.wg.Add(1)
		go func() {
			tuiApp.RerenderCreds()
			tuiApp.wg.Done()
		}()

	})

	loginForm.AddButton("Cancel", func() {
		f := tuiApp.FormLogin.GetItem(0).(*tview.Form)
		f = clearForm(f)
		tuiApp.Pages.SwitchToPage(PageMain)
		tuiApp.App.SetFocus(tuiApp.FormAuth)
	})

	loginForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			tuiApp.log.Info().
				Msg("caught esc button")

			f := tuiApp.FormLogin.GetItem(0).(*tview.Form)
			f = clearForm(f)
			tuiApp.Pages.SwitchToPage(PageMain)
			tuiApp.App.SetFocus(tuiApp.FormAuth)
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

		if err := tuiApp.Usecase.Register(tuiApp.Ctx, email, username, password, repeat); err != nil {
			stErr := status.FromContextError(err)
			if strings.Contains(stErr.Message(), "connection refused") {
				tuiApp.log.Error().Err(err).Msg("server is unavailable")
				errAuthForm := NewModalWithParams(tuiApp, "Server is unavailable!", PageRegister)
				tuiApp.Pages.AddPage(PageRegisterError, errAuthForm, true, false)
				tuiApp.Pages.SwitchToPage(PageRegisterError)
				return
			}

			tuiApp.log.Error().Err(err).Msg("failed to register")
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
