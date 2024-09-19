package tui

import (
	"strconv"
	"strings"

	"github.com/rivo/tview"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
)

var (
	HintTextCards = genHelp("card")
)

type CardDetails struct {
	*tview.Form

	//tuiApp *TUI

	FieldName       *tview.InputField
	FieldBank       *tview.DropDown
	FieldPerson     *tview.InputField
	FieldNumber     *tview.InputField
	FieldExpiration *tview.InputField
	FieldCVC        *tview.InputField
	FieldPIN        *tview.InputField
	FieldDesc       *tview.TextArea

	BtnSave   *tview.Button
	BtnCancel *tview.Button

	SaveLabel   string
	CancelLabel string

	CurrentCard *entities.Card

	// fields sizes
	FieldWidth  int
	FieldHeight int
	maxSigns    int
}

func NewCardDetails(card *entities.Card) *CardDetails {
	if card == nil {
		card = &entities.Card{}
	}

	bank := tview.NewDropDown().
		SetLabel("Bank:").
		SetOptions(entities.Banks, nil)
	bank.SetTitle(card.Bank)
	bank.SetTextOptions(card.Bank, "", "", "", card.Bank)

	form := &CardDetails{
		Form:            tview.NewForm(),
		FieldName:       tview.NewInputField().SetLabel("Name:").SetText(card.Name),
		FieldBank:       bank,
		FieldPerson:     tview.NewInputField().SetLabel("Person:").SetText(card.Person),
		FieldNumber:     tview.NewInputField().SetLabel("Number:").SetText(strconv.Itoa(card.Number)),
		FieldExpiration: tview.NewInputField().SetLabel("Expiration").SetText(card.Expiration),
		FieldCVC:        tview.NewInputField().SetLabel("CVC").SetText(strconv.Itoa(card.CVC)),
		FieldPIN:        tview.NewInputField().SetLabel("PIN").SetText(strconv.Itoa(card.PIN)),
		FieldDesc:       tview.NewTextArea().SetLabel("Description:").SetText(card.Description, true),
		CurrentCard:     card,
	}

	form.Form.SetBorder(true).
		SetTitle("Details")

	form.Form.AddFormItem(form.FieldName)
	form.Form.AddFormItem(form.FieldBank)
	form.Form.AddFormItem(form.FieldPerson)
	form.Form.AddFormItem(form.FieldNumber)
	form.Form.AddFormItem(form.FieldExpiration)
	form.Form.AddFormItem(form.FieldCVC)
	form.Form.AddFormItem(form.FieldPIN)
	form.Form.AddFormItem(form.FieldDesc)

	return form
}

func (form *CardDetails) Rerender(card *entities.Card) {
	form.FieldName.SetText(card.Name)
	form.FieldBank.SetTextOptions(card.Bank, "", "", "", card.Bank)
	form.FieldPerson.SetText(card.Person)
	form.FieldNumber.SetText(strconv.Itoa(card.Number))
	form.FieldExpiration.SetText(card.Expiration)
	form.FieldCVC.SetText(strconv.Itoa(card.CVC))
	form.FieldPIN.SetText(strconv.Itoa(card.PIN))
	form.FieldDesc.SetText(card.Description, true)
	form.CurrentCard = card
}

func (form *CardDetails) Add(tuiApp *TUI, ind int, list *CardList) {
	form.EmptyFields()

	form.AddButton("Save", func() {
		//defer continue cred sync
		defer tuiApp.Usecase.ContinueSync()

		// defer remove buttons
		defer form.HideButtons()

		newCard, err := form.GetCurrentValues()
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add credential on server side")
			errModal := NewModalWithParams(tuiApp, err.Error(), PageCredsMenu)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		newCard.Type = entities.UserCard
		newCard.ID = hashes.GeneratePassID2()

		if err := tuiApp.Usecase.AddBlob(tuiApp.Ctx, newCard); err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add credential on server side")
			errModal := NewModalWithParams(tuiApp, err.Error(), PageCredsMenu)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		form.Rerender(newCard)

		// rerender credsList
		tuiApp.App.SetFocus(list)

	})

	form.AddButton("Cancel", func() {
		//defer continue cred sync
		defer tuiApp.Usecase.ContinueSync()

		// defer remove buttons
		defer form.HideButtons()
		// rerender credsList
		defer tuiApp.App.SetFocus(list)
		if tuiApp.Usecase.CredsLen() > 0 {
			card, _ := tuiApp.Usecase.GetCardByIND(ind)
			form.Rerender(card)
			return
		}

		// clear fields if there isn't any blobsUC
		//form.EmptyFields()
		form.HideFields()
	})

}

func (form *CardDetails) HideButtons() {
	if form.GetButtonCount() <= 0 {
		return
	}

	form.RemoveButton(1)
	form.RemoveButton(0)
}

func (form *CardDetails) HideFields() {
	for ind := form.Form.GetFormItemCount() - 1; ind >= 0; ind-- {
		form.Form.RemoveFormItem(ind)
	}
}

func (form *CardDetails) EmptyFields() {
	if form.FieldName != nil {
		form.FieldName.SetText("")
	}

	if form.FieldBank != nil {
		form.FieldBank.SetTextOptions("", "", "", "", "")
	}

	if form.FieldPerson != nil {
		form.FieldPerson.SetText("")
	}

	if form.FieldNumber != nil {
		form.FieldNumber.SetText("")
	}

	if form.FieldExpiration != nil {
		form.FieldExpiration.SetText("")
	}

	if form.FieldCVC != nil {
		form.FieldCVC.SetText("")
	}

	if form.FieldPIN != nil {
		form.FieldPIN.SetText("")
	}

	if form.FieldDesc != nil {
		form.FieldDesc.SetText("", true)
	}
}
func (form *CardDetails) GetCurrentValues() (newCard *entities.Card, err error) {
	newCard = new(entities.Card)

	newCard.Name = form.FieldName.GetText()
	_, newCard.Bank = form.FieldBank.GetCurrentOption()
	newCard.Person = form.FieldPerson.GetText()
	newCard.Description = form.FieldDesc.GetText()

	numberStr := form.FieldNumber.GetText()
	numberStr = strings.ReplaceAll(numberStr, " ", "")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, err
	}
	newCard.Number = number

	newCard.Expiration = form.FieldExpiration.GetText()

	cvc, err := strconv.Atoi(form.FieldCVC.GetText())
	if err != nil {
		return nil, err
	}
	newCard.CVC = cvc

	pin, err := strconv.Atoi(form.FieldPIN.GetText())
	if err != nil {
		return nil, err
	}
	newCard.PIN = pin

	return newCard, nil
}
