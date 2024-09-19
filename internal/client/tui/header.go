package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Header struct {
	*tview.Flex
	curID int
}

func NewHeader(ind int) *Header {
	ind += 1 // to skip the first empty item

	pgView := tview.NewFlex().
		AddItem(tview.NewTextView(), 0, 1, false).
		AddItem(tview.NewTextView().SetText("F2 - Resources"), 0, 3, false).
		AddItem(tview.NewTextView().SetText("F3 - Bank cards"), 0, 3, false).
		AddItem(tview.NewTextView().SetText("F4 - User Notes"), 0, 3, false).
		AddItem(tview.NewTextView().SetText("F5 - Files"), 0, 3, false).
		AddItem(tview.NewTextView(), 0, 1, false)
	pgView.SetBorder(true).
		SetBorderPadding(0, 0, 0, 0)

	if ind < pgView.GetItemCount() {
		item := pgView.GetItem(ind)
		view := item.(*tview.TextView)
		view.SetTextColor(tcell.ColorYellow)
		text := view.GetText(true)
		view.SetText("*" + text)
	}

	return &Header{
		Flex:  pgView,
		curID: ind,
	}

}

func (h *Header) ChangePage(ind int) {
	// already switched, do nothing
	if ind == h.curID {
		return
	}

	cnt := h.Flex.GetItemCount()

	// restore current to defaults
	if h.curID < cnt {
		item := h.Flex.GetItem(h.curID)
		view := item.(*tview.TextView)
		view.SetTextColor(tcell.ColorWhite)
		text := view.GetText(true)
		view.SetText(strings.Trim(text, "*"))
	}

	// set new ID as current
	if ind < cnt {
		item := h.Flex.GetItem(ind)
		view := item.(*tview.TextView)
		view.SetTextColor(tcell.ColorYellow)
		text := view.GetText(true)
		view.SetText("*" + text)
	}

	h.curID = ind
}
