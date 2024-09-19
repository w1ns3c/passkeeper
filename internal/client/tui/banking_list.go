package tui

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"passkeeper/internal/config"
	"passkeeper/internal/entities"
)

func NewBanking(cards []*entities.Card) *tview.Flex {
	var (
		form *CardDetails
	)

	list := NewCardList(cards)
	list.Rerender(cards)

	ind := list.GetCurrentItem()
	if cards != nil && len(cards) > ind {
		form = NewCardDetails(cards[ind])
	} else {
		form = NewCardDetails(nil)
	}

	list.SetBorder(true).
		SetTitle("Bank Cards")

	helpCards := tview.NewTextView().
		SetTextColor(tcell.ColorBisque).
		SetText(HintTextCards)

	subFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 2, false).
		AddItem(helpCards, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(list, 0, 2, true).
		AddItem(subFlex, 0, 3, false)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//ind := list.GetCurrentItem()
		delFunc := func() {
			if list.GetItemCount() == 0 {
				return
			}

			//delConfirm := DeleteModal(tuiApp, ind)
			//pageConfirm := "confirmation"
			//tuiApp.Pages.AddPage(pageConfirm, delConfirm, true, false)
			//tuiApp.Pages.SwitchToPage(pageConfirm)
		}

		// inputs for credList
		switch event.Key() {
		case tcell.KeyEsc:
			fallthrough
		case tcell.KeyCtrlL:

			//pageName := "Logout"
			//logoutModal := LogoutModal(tuiApp)
			//tuiApp.Pages.AddPage(pageName, logoutModal, true, false)
			//tuiApp.Pages.SwitchToPage(pageName)

		case tcell.KeyDelete:
			delFunc()
		}

		switch event.Rune() {
		case 'a':
			//tuiApp.Usecase.StopSync()
			//tuiApp.App.SetFocus(viewForm)
			//viewForm.Add(ind, credList)
		case 'e':
			//tuiApp.Usecase.StopSync()
			//tuiApp.App.SetFocus(viewForm)
			//viewForm.Edit(ind, credList)

		case 'd':
			delFunc()
		case ' ':
			//showFunc(false)
			//viewForm.ShowSwitchPass()
		}

		return event
	})

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		curID := list.GetCurrentItem()
		form.Rerender(cards[curID])
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		form.BtnSave = tview.NewButton("Save")
		form.BtnCancel = tview.NewButton("Cancel")

		form.BtnSave.SetSelectedFunc(func() {

		})

		return event
	})

	return flex
}

type CardList struct {
	*tview.List
	cards []*entities.Card
}

func NewCardList(cards []*entities.Card) *CardList {
	list := tview.NewList()
	list.ShowSecondaryText(false).
		SetBorderPadding(0, 0, 0, 0)

	return &CardList{
		List:  list,
		cards: cards,
	}
}

func (list *CardList) Rerender(cards []*entities.Card) {
	for ind := list.GetItemCount() - 1; ind >= 0; ind-- {
		list.RemoveItem(ind)
	}

	if cards != nil {
		for ind, card := range cards {
			res := GenCardShortName(card)
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

func GenCardShortName(card *entities.Card) string {
	if card == nil {
		return ""
	}

	var (
		res string
		m   = config.MaxNameLen
	)

	if len(card.Name) > m {
		parts := strings.Split(card.Name, " ")
		if len(parts) == 1 {
			res = parts[0]
		} else {
			res = strings.Join(parts[:2], "_")[:m]
		}
	} else {
		res = card.Name
	}

	if card.Number != 0 {
		last4sign := card.Number % 10000
		res += " - *" + strconv.Itoa(last4sign)
	}

	return res
}
