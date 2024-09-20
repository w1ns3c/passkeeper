package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
)

var (
	HintTextCards           = genHelp("card")
	bankingExpirationFormat = "01/06"
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
		Form:        tview.NewForm(),
		FieldName:   tview.NewInputField().SetLabel("Name:").SetText(card.Name),
		FieldBank:   bank,
		FieldPerson: tview.NewInputField().SetLabel("Person:").SetText(card.Person),
		FieldNumber: tview.NewInputField().SetLabel("Number:").SetText(strconv.Itoa(card.Number)),
		FieldExpiration: tview.NewInputField().SetLabel("Expiration").
			SetText(card.Expiration.Format(bankingExpirationFormat)),
		FieldCVC:    tview.NewInputField().SetLabel("CVC").SetText(strconv.Itoa(card.CVC)),
		FieldPIN:    tview.NewInputField().SetLabel("PIN").SetText(strconv.Itoa(card.PIN)),
		FieldDesc:   tview.NewTextArea().SetLabel("Description:").SetText(card.Description, true),
		CurrentCard: card,
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
	form.FieldBank.SetTextOptions("", "", "", "", card.Bank)
	form.FieldPerson.SetText(card.Person)

	num := strconv.Itoa(card.Number)
	if num == "0" {
		form.FieldNumber.SetText("")
	} else {
		form.FieldNumber.SetText(num)
	}

	form.FieldExpiration.SetText(card.Expiration.Format(bankingExpirationFormat))

	cvc := strconv.Itoa(card.CVC)
	if cvc == "0" {
		form.FieldCVC.SetText("")
	} else {
		form.FieldCVC.SetText(cvc)
	}

	pin := strconv.Itoa(card.PIN)
	if pin == "0" {
		form.FieldPIN.SetText("")
	} else {
		form.FieldPIN.SetText(pin)
	}

	form.FieldDesc.SetText(card.Description, true)
	form.CurrentCard = card
}

func (form *CardDetails) Add(tuiApp *TUI, ind int, list *CardList) {
	form.EmptyFields()

	form.AddButton("Save", func() {
		newCard, err := form.GetCurrentValues()
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add new card on client side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		newCard.Type = entities.UserCard
		newCard.ID = hashes.GeneratePassID2()

		if err := tuiApp.Usecase.AddBlob(tuiApp.Ctx, newCard); err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add new card on server side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		// rerender credsList
		list.Rerender(tuiApp.Usecase.GetCards())
		form.Rerender(newCard)
		tuiApp.App.SetFocus(list)

		// hide buttons
		form.HideButtons()

		//defer continue cred sync
		tuiApp.Usecase.ContinueSync()

		tuiApp.log.Info().
			Str("id", newCard.ID).Msg("cred added")
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

		curCard, err := tuiApp.Usecase.GetCardByIND(list.GetCurrentItem())
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("can't get current card")

		}

		form.Rerender(curCard)
	})

}

func (form *CardDetails) Edit(tuiApp *TUI, ind int, list *CardList) {
	tmp := tuiApp.Usecase.GetCards()
	if tmp == nil || len(tmp) <= ind || len(tmp) == 0 {
		return
	}

	form.AddButton("Save", func() {
		cur, err := tuiApp.Usecase.GetCardByIND(ind)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to edit new card on client side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)

			return
		}

		editedCard, err := form.GetCurrentValues()
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to edit new card on client side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)

			return
		}

		editedCard.SetID(cur.GetID())

		if err := tuiApp.Usecase.EditBlob(tuiApp.Ctx, editedCard, ind); err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to edit new card on server side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)

			return
		}

		// rerender credsList
		list.Rerender(tuiApp.Usecase.GetCards())
		form.Rerender(editedCard)
		tuiApp.App.SetFocus(list)

		// hide buttons
		form.HideButtons()

		//defer continue cred sync
		tuiApp.Usecase.ContinueSync()

		tuiApp.log.Info().
			Str("id", editedCard.ID).Msg("cred updated")
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

// HideButtons hide Save/Cancel buttons
func (form *CardDetails) HideButtons() {
	if form.GetButtonCount() <= 0 {
		return
	}

	form.RemoveButton(1)
	form.RemoveButton(0)
}

// HideFields remove all items from form
func (form *CardDetails) HideFields() {
	for ind := form.Form.GetFormItemCount() - 1; ind >= 0; ind-- {
		form.Form.RemoveFormItem(ind)
	}
}

// EmptyFields reset fields' values
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

// GetCurrentValues get values from user input and format it to Card entity
func (form *CardDetails) GetCurrentValues() (newCard *entities.Card, err error) {
	newCard = new(entities.Card)
	newCard.Type = entities.UserCard

	newCard.Name = form.FieldName.GetText()
	_, newCard.Bank = form.FieldBank.GetCurrentOption()
	newCard.Person = form.FieldPerson.GetText()
	newCard.Description = form.FieldDesc.GetText()

	numberStr := form.FieldNumber.GetText()
	numberStr = strings.ReplaceAll(numberStr, " ", "")
	number, err := strconv.Atoi(numberStr)
	if err != nil {

		return nil, fmt.Errorf("can't parse card number")
	}
	newCard.Number = number

	exp := form.FieldExpiration.GetText()
	newCard.Expiration, err = time.Parse(bankingExpirationFormat, exp)
	if err != nil {
		return nil, fmt.Errorf("can't parse card expiration")
	}

	cvc, err := strconv.Atoi(form.FieldCVC.GetText())
	if err != nil {

		return nil, fmt.Errorf("can't parse card CVV")
	}
	newCard.CVC = cvc

	pin, err := strconv.Atoi(form.FieldPIN.GetText())
	if err != nil {

		return nil, fmt.Errorf("can't parse card CVV")
	}
	newCard.PIN = pin

	return newCard, nil
}
