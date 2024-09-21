package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"passkeeper/internal/entities"
)

var (
	hintTextCreds = genHelp("card")
	hintCreds     = tview.NewTextView().
			SetTextColor(tcell.ColorBisque).
			SetText(hintTextCreds)

	passHidden = "******"
)

type CredListInf interface {
	tview.Primitive
	Rerender()
	SetItemText(index int, main string, secondary string)
}

type CredList struct {
	*tview.List
	tuiApp *TUI
}

func (list *CredList) Rerender() {
	for ind := list.GetItemCount() - 1; ind >= 0; ind-- {
		list.RemoveItem(ind)
	}

	// fill list fields with values
	creds := list.tuiApp.Usecase.GetCreds()
	if creds != nil {
		for ind, cred := range creds {
			res := FilterResource(cred.Resource)
			if ind < 9 {
				list.AddItem(res, "", rune(49+ind), nil)
			} else if ind == 9 {
				list.AddItem(res, "", 'X', nil)
			} else {
				list.AddItem(res, "", rune(65+ind-10), nil)

			}
		}
	}
}

func (list *CredList) SetItemText(index int, main string, secondary string) {
	list.List.SetItemText(index, main, secondary)
}

func NewList(tuiApp *TUI) *CredList {
	return &CredList{
		List: tview.NewList().
			ShowSecondaryText(false),
		tuiApp: tuiApp,
	}
}

func NewCredsList(tuiApp *TUI) *tview.Flex {

	credList := NewList(tuiApp)
	credList.Rerender()

	credList.SetBorder(true).
		SetBorderPadding(0, 0, 0, 0).
		SetTitle("Credentials")

	var viewForm = NewDetailsForm(tuiApp)
	if tuiApp.Usecase.CredsNotNil() {
		if tuiApp.Usecase.CredsLen() != 0 {
			cred, _ := tuiApp.Usecase.GetCredByIND(0)
			viewForm.ShowFields()
			viewForm.EmptyFields()
			viewForm.SetHiddenCred(cred)
		}
	}

	//Hints in full length
	listFlex := tview.NewFlex().
		AddItem(credList, 0, 2, true).
		AddItem(viewForm, 0, 3, true)

	fullFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		//AddItem(NewHeader(1), 0, 1, false).
		AddItem(listFlex, 0, 10, true).
		AddItem(hintCreds, 0, 2, false)

	credList.SetChangedFunc(func(ind int, mainText string, secondaryText string, shortcut rune) {
		viewForm.ShowFields()
		viewForm.EmptyFields()
		cred, err := tuiApp.Usecase.GetCredByIND(ind)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("wrong ind")
			return
		}

		viewForm.SetCurrentCred(cred)
		viewForm.HidePassword()
	})

	credList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ind := credList.GetCurrentItem()
		delFunc := func() {
			if credList.GetItemCount() == 0 {
				return
			}

			delConfirm := DeleteModal(tuiApp, ind, entities.BlobCred)
			pageConfirm := "confirmation"
			tuiApp.Pages.AddPage(pageConfirm, delConfirm, true, false)
			tuiApp.Pages.SwitchToPage(pageConfirm)
		}

		// inputs for credList
		switch event.Key() {
		case tcell.KeyEsc:
			fallthrough
		case tcell.KeyCtrlL:

			pageName := "Logout"
			logoutModal := LogoutModal(tuiApp)
			tuiApp.Pages.AddPage(pageName, logoutModal, true, false)
			tuiApp.Pages.SwitchToPage(pageName)

		case tcell.KeyDelete:
			delFunc()
		case tcell.KeyEnter:
			viewForm.ShowSwitchPass()
		}

		switch event.Rune() {
		case 'a':
			tuiApp.Usecase.StopSync()
			tuiApp.App.SetFocus(viewForm)
			viewForm.Add(ind, credList)
		case 'e':
			tuiApp.Usecase.StopSync()
			tuiApp.App.SetFocus(viewForm)
			viewForm.Edit(ind, credList)

		case 'd':
			delFunc()
		case ' ':
			//showFunc(false)
			viewForm.ShowSwitchPass()
		}
		return event
	})

	viewForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			fallthrough
		case tcell.KeyEscape:
			tuiApp.Usecase.ContinueSync()
			ind := credList.GetCurrentItem()
			viewForm.HideButtons()
			viewForm.EmptyFields()
			viewForm.tuiApp.App.SetFocus(credList)
			l := tuiApp.Usecase.CredsLen()
			cred, _ := tuiApp.Usecase.GetCredByIND(ind)
			//if tuiApp.Creds != nil {
			if tuiApp.Usecase.CredsNotNil() {
				if l > 0 {
					if l > ind {
						viewForm.SetCurrentCred(cred)
						viewForm.Rerender()
						return event
					}
				}
			}
			viewForm.HideFields()
		}
		return event
	})

	return fullFlex
}
