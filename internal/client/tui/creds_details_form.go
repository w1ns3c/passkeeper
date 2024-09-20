package tui

import (
	"time"

	"github.com/rivo/tview"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
)

type Details struct {
	*tview.Form

	tuiApp *TUI

	FieldDate  *tview.InputField
	FieldRes   *tview.InputField
	FieldLogin *tview.InputField
	FieldPass  *tview.InputField
	FieldDesc  *tview.TextArea

	BtnSave   *tview.Button
	BtnCancel *tview.Button

	SaveLabel   string
	CancelLabel string

	CurrentCred  *entities.Credential
	HiddenPass   bool
	CurrentField int

	// fields sizes
	FieldWidth  int
	FieldHeight int
	maxSigns    int

	PassValue string
}

func NewDetailsForm(tuiApp *TUI) *Details {
	form := &Details{
		tuiApp:      tuiApp,
		Form:        tview.NewForm(),
		HiddenPass:  false,
		SaveLabel:   "Save",
		CancelLabel: "Cancel",
		FieldWidth:  40,
		FieldHeight: 6,
	}
	form.
		SetBorder(true).
		SetTitle("Details")

	form.FieldDate = tview.NewInputField().
		SetLabel("Date").
		SetFieldWidth(form.FieldWidth)

	form.FieldDate.
		SetDisabled(true)

	form.FieldRes = tview.NewInputField().
		SetLabel("Resource").
		SetFieldWidth(form.FieldWidth)

	form.FieldLogin = tview.NewInputField().
		SetLabel("Login").
		SetFieldWidth(form.FieldWidth)

	form.FieldPass = tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(form.FieldWidth)

	form.FieldDesc = tview.NewTextArea().
		SetLabel("Description").
		SetSize(form.FieldHeight, form.FieldWidth).
		SetMaxLength(form.maxSigns)

	return form
}

func (form *Details) HideButtons() {
	if form.GetButtonCount() <= 0 {
		return
	}

	form.RemoveButton(1)
	form.RemoveButton(0)
}

// Add switch tui app to cred adding form
func (form *Details) Add(ind int, list CredListInf) {
	//form.ShowItems()
	form.HideFields()
	form.ShowPartFields()
	form.EmptyFields()
	form.tuiApp.App.SetFocus(form)

	form.AddButton("Save", func() {
		//defer continue cred sync
		defer form.tuiApp.Usecase.ContinueSync()

		// defer remove buttons
		defer form.HideButtons()

		res, login, password, desc := form.GetCurrentValues()

		newCred := &entities.Credential{
			Type:        entities.UserCred,
			ID:          hashes.GeneratePassID2(),
			Date:        time.Now(),
			Resource:    res,
			Login:       login,
			Password:    password,
			Description: desc,
		}

		if err := form.tuiApp.Usecase.AddBlob(form.tuiApp.Ctx, newCred); err != nil {
			form.tuiApp.log.Error().
				Err(err).Msg("failed to add credential on server side")
			errModal := NewModalWithParams(form.tuiApp, err.Error(), SubPageCreds)
			form.tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			form.tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		form.Rerender()
		form.ShowPassword()

		// rerender credsList
		form.tuiApp.App.SetFocus(list)

		list.Rerender()
	})

	form.AddButton("Cancel", func() {
		//defer continue cred sync
		defer form.tuiApp.Usecase.ContinueSync()

		// defer remove buttons
		defer form.cancel()
		// rerender credsList
		defer form.tuiApp.App.SetFocus(list)
		if form.tuiApp.Usecase.CredsLen() > 0 {
			form.Rerender()
			cred, _ := form.tuiApp.Usecase.GetCredByIND(ind)
			form.SetCurrentCred(cred)
			return
		}

		// clear fields if there isn't any blobsUC
		form.EmptyFields()
		form.HideFields()
	})
}

// Edit switch tui app to cred editing form
func (form *Details) Edit(ind int, list CredListInf) {
	// we shouldn't edit unexisted credential
	l := form.tuiApp.Usecase.CredsLen()
	if l <= ind || l == 0 {
		return
	}

	form.ShowItems()
	form.tuiApp.App.SetFocus(form.FieldRes)

	form.AddButton("Save", func() {
		//defer continue cred sync
		defer form.tuiApp.Usecase.ContinueSync()

		// defer remove buttons
		defer form.HideButtons()

		res, login, password, desc := form.GetCurrentValues()

		cred, err := form.tuiApp.Usecase.GetCredByIND(ind)
		if err != nil {
			form.tuiApp.log.Error().
				Err(err).Msg("failed to edit credential on client side")
			errModal := NewModalWithParams(form.tuiApp, err.Error(), SubPageCreds)
			form.tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			form.tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		cred.Date = time.Now()
		cred.Resource = res
		cred.Login = login
		cred.Password = password
		cred.Description = desc

		// send creds to server
		if err := form.tuiApp.Usecase.EditBlob(form.tuiApp.Ctx, cred, ind); err != nil {
			form.tuiApp.log.Error().
				Err(err).Msg("failed to edit credential on server side")
			errModal := NewModalWithParams(form.tuiApp, err.Error(), SubPageCreds)
			form.tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			form.tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		//list.SetItemText(ind, resFilter, "")
		resFilter := FilterResource(res)
		list.SetItemText(ind, resFilter, "")
		form.ShowPassword()
		form.tuiApp.App.SetFocus(list)
	})

	form.AddButton("Cancel", func() {
		//defer continue cred sync
		defer form.tuiApp.Usecase.ContinueSync()

		defer form.HideButtons() // remove buttons from form
		form.Rerender()
		cred, _ := form.tuiApp.Usecase.GetCredByIND(ind)
		form.SetCurrentCred(cred)
		form.tuiApp.App.SetFocus(list)
	})

}

func (form *Details) cancel() {
	defer form.HideButtons()  // defer remove buttons
	defer form.HidePassword() // hide password
}

func (form *Details) Rerender() {
	form.HideFields()
	form.ShowFields()
}

func (form *Details) ShowFields() {
	if form.GetFormItemCount() > 0 {
		return
	}

	form.Form.AddFormItem(form.FieldRes)
	form.Form.AddFormItem(form.FieldLogin)
	form.Form.AddFormItem(form.FieldPass)
	form.Form.AddFormItem(form.FieldDate)
	form.Form.AddFormItem(form.FieldDesc)

	//form.CurrentField = 0
	//form.SetFocus(form.CurrentField)

}

func (form *Details) ShowPartFields() {
	if form.GetFormItemCount() > 0 {
		return
	}

	form.Form.AddFormItem(form.FieldRes)
	form.Form.AddFormItem(form.FieldLogin)
	form.Form.AddFormItem(form.FieldPass)
	//form.Form.AddFormItem(form.FieldDate)
	form.Form.AddFormItem(form.FieldDesc)

	form.CurrentField = 0
	form.SetFocus(form.CurrentField)

}

func (form *Details) ShowItems() {
	form.ShowFields()
	form.ShowPassword()
}

func (form *Details) keyUp() {
	minInd := 0
	maxInd := form.GetFormItemCount() - 1

	if maxInd < 0 {
		return
	}

	if form.CurrentField == minInd {
		form.CurrentField = maxInd
		form.SetFocus(form.CurrentField)
		return
	}

	form.CurrentField--
	form.SetFocus(form.CurrentField)

}

func (form *Details) keyDown() {

	maxInd := form.GetFormItemCount() - 1

	if maxInd < 0 {
		return
	}

	if form.CurrentField == maxInd {
		// switch to the first item
		form.CurrentField = 0
		form.SetFocus(form.CurrentField)
		return
	}

	form.CurrentField++
	form.SetFocus(form.CurrentField)

}

func (form *Details) HideFields() {
	for ind := form.Form.GetFormItemCount() - 1; ind >= 0; ind-- {
		form.Form.RemoveFormItem(ind)
	}

}

// ResetFields is alias for EmptyFields
func (form *Details) ResetFields() {
	form.EmptyFields()
}

// EmptyFields method clear input fields
func (form *Details) EmptyFields() {
	if form.FieldRes != nil {
		form.FieldRes.SetText("")
	}

	if form.FieldDate != nil {
		form.FieldDate.SetText("")
	}

	if form.FieldLogin != nil {
		form.FieldLogin.SetText("")
	}

	if form.FieldPass != nil {
		form.FieldPass.SetText("")
	}

	if form.FieldDesc != nil {
		form.FieldDesc.SetText("", true)
	}
}

func (form *Details) SetHiddenCred(cred *entities.Credential) *Details {
	if form.FieldRes != nil {
		form.FieldRes.SetText(cred.Resource)
	}
	if form.FieldDate != nil {
		form.FieldDate.SetText(cred.Date.Format("2006.01.02 15:04:05"))
	}
	if form.FieldLogin != nil {
		form.FieldLogin.SetText(cred.Login)
	}
	if form.FieldPass != nil {
		form.PassValue = cred.Password
		form.FieldPass.SetText(PassHidden)
		form.HiddenPass = true
	}
	if form.FieldDesc != nil {
		form.FieldDesc.SetText(cred.Description, true)
	}
	return form
}

func (form *Details) SetCurrentCred(cred *entities.Credential) *Details {
	if form.FieldRes != nil {
		form.FieldRes.SetText(cred.Resource)
	}
	if form.FieldDate != nil {
		form.FieldDate.SetText(cred.Date.Format("2006.01.02 15:04:05"))
	}
	if form.FieldLogin != nil {
		form.FieldLogin.SetText(cred.Login)
	}
	if form.FieldPass != nil {
		form.PassValue = cred.Password
		form.FieldPass.SetText(form.PassValue)
	}
	if form.FieldDesc != nil {
		form.FieldDesc.SetText(cred.Description, true)
	}
	return form
}

func (form *Details) HidePassword() *Details {
	if form.FieldPass == nil {
		return form
	}

	form.tuiApp.log.Info().Msg("hide current cred password")
	form.HiddenPass = true
	form.PassValue = form.FieldPass.GetText()
	if form.PassValue != "" {
		form.FieldPass.SetText(PassHidden)
	}
	return form
}

func (form *Details) ShowPassword() {
	form.tuiApp.log.Info().Msg("show current cred password")
	form.HiddenPass = false
	form.FieldPass.SetText(form.PassValue)
}

func (form *Details) ShowSwitchPass() {
	if form.HiddenPass {
		form.ShowPassword()
	} else {
		form.HidePassword()
	}
}

// FillFields alias for SetCurrentCred
func (form *Details) FillFields(cred *entities.Credential) {
	form.SetCurrentCred(cred)
}

func (form *Details) GetCurrentValues() (resource, login, password, description string) {
	resource = form.FieldRes.GetText()
	login = form.FieldLogin.GetText()
	form.PassValue = form.FieldPass.GetText()
	password = form.PassValue
	description = form.FieldDesc.GetText()

	return resource, login, password, description
}

func (form *Details) FieldsIsEmpty() bool {
	return form.FieldRes.GetText() == "" &&
		form.FieldLogin.GetText() == "" &&
		form.FieldPass.GetText() == "" &&
		form.FieldDesc.GetText() == ""
}
