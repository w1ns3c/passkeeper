package tui

import (
	"strconv"

	"github.com/rivo/tview"

	"passkeeper/internal/entities"
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
	bank.AddOption(card.Bank, nil)

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
	form.FieldBank.SetTitle(card.Bank)
	form.FieldPerson.SetText(card.Person)
	form.FieldNumber.SetText(strconv.Itoa(card.Number))
	form.FieldExpiration.SetText(card.Expiration)
	form.FieldCVC.SetText(strconv.Itoa(card.CVC))
	form.FieldPIN.SetText(strconv.Itoa(card.PIN))
	form.FieldDesc.SetText(card.Description, true)
}
