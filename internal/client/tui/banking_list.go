package tui

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/structs"
)

func (tuiApp *TUI) NewBanking(cards []*structs.Card) *tview.Flex {
	var (
		viewForm *CardDetails
	)

	list := NewCardList(cards)
	list.Rerender(cards)

	ind := list.GetCurrentItem()
	if cards != nil && len(cards) > ind {
		viewForm = NewCardDetails(cards[ind])
	} else {
		viewForm = NewCardDetails(nil)
	}

	list.SetBorder(true).
		SetTitle("Bank Cards")

	helpCards := tview.NewTextView().
		SetTextColor(tcell.ColorBisque).
		SetText(hintTextCards)

	subFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(viewForm, 0, 2, false).
		AddItem(helpCards, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(list, 0, 2, true).
		AddItem(subFlex, 0, 3, false)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// inputs for bank cards
		switch event.Key() {
		case tcell.KeyEsc:
			fallthrough
		case tcell.KeyCtrlL:

			pageName := "Logout"
			logoutModal := LogoutModal(tuiApp)
			tuiApp.Pages.AddPage(pageName, logoutModal, true, false)
			tuiApp.Pages.SwitchToPage(pageName)

		case tcell.KeyDelete:
			list.Delete(tuiApp, ind)
		}

		switch event.Rune() {
		case 'a':
			tuiApp.Usecase.StopSync()
			tuiApp.App.SetFocus(viewForm)
			viewForm.Add(tuiApp, ind, list)
		case 'e':
			tuiApp.Usecase.StopSync()
			tuiApp.App.SetFocus(viewForm)
			viewForm.Edit(tuiApp, ind, list)

		case 'd':
			list.Delete(tuiApp, ind)
		}

		return event
	})

	list.SetChangedFunc(func(ind int, mainText string, secondaryText string, shortcut rune) {
		card, err := tuiApp.Usecase.GetCardByIND(ind)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("wrong card ind")
			return
		}
		viewForm.Rerender(card)
	})

	return flex
}

type CardList struct {
	*tview.List
	cards []*structs.Card
}

func NewCardList(cards []*structs.Card) *CardList {
	list := tview.NewList()
	list.ShowSecondaryText(false).
		SetBorderPadding(0, 0, 0, 0)

	return &CardList{
		List:  list,
		cards: cards,
	}
}

func (list *CardList) Rerender(cards []*structs.Card) {
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

// GenCardShortName beautify card name to show it in the list
func GenCardShortName(card *structs.Card) string {
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

func (list *CardList) Delete(tuiApp *TUI, ind int) {
	if list.GetItemCount() == 0 {
		return
	}

	delConfirm := DeleteModal(tuiApp, ind, structs.BlobCard)
	pageConfirm := "confirmation"
	tuiApp.Pages.AddPage(pageConfirm, delConfirm, true, false)
	tuiApp.Pages.SwitchToPage(pageConfirm)
}
