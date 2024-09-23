package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/structs"
)

var (
	HintTextNotes = "note"
)

func (tuiApp *TUI) NewNotes(notes []*structs.Note) *tview.Flex {
	var viewForm *NoteDetails

	list := NewNotesList(notes)
	list.Rerender(notes)

	ind := list.GetCurrentItem()
	if notes != nil && len(notes) > ind {
		viewForm = NewNoteDetails(notes[ind])
	} else {
		viewForm = NewNoteDetails(nil)
	}

	list.SetBorder(true).
		SetTitle("Notes")

	helpNotes := tview.NewTextView().
		SetTextColor(tcell.ColorBisque).
		SetText(genHelp(HintTextNotes))

	subFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(viewForm, 0, 2, false).
		AddItem(helpNotes, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(list, 0, 2, true).
		AddItem(subFlex, 0, 3, false)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

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
		note, err := tuiApp.Usecase.GetNoteByIND(ind)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("wrong note ind")

			return
		}
		viewForm.Rerender(note)
	})

	return flex
}

type NotesList struct {
	*tview.List
	notes []*structs.Note
}

func NewNotesList(notes []*structs.Note) *NotesList {
	list := tview.NewList()
	list.ShowSecondaryText(false).
		SetBorderPadding(0, 0, 0, 0)
	return &NotesList{
		List:  list,
		notes: notes,
	}
}

func (list *NotesList) Rerender(notes []*structs.Note) {
	for ind := list.GetItemCount() - 1; ind >= 0; ind-- {
		list.RemoveItem(ind)
	}

	if notes != nil {
		for ind, note := range notes {
			res := GenNoteShortName(note)
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

func GenNoteShortName(note *structs.Note) string {
	if note == nil {
		return ""
	}

	var (
		res string
		m   = config.MaxNameLen
	)

	if len(note.Name) > m {
		parts := strings.Split(note.Name, " ")
		if len(parts) == 1 {
			res = parts[0]
		} else {
			res = strings.Join(parts[:1], "_")
			if len(res) > m {
				res = res[:m]
			}
		}
	} else {
		res = note.Name
	}

	return res
}

func (list *NotesList) Delete(tuiApp *TUI, ind int) {
	if list.GetItemCount() == 0 {
		return
	}

	delConfirm := DeleteModal(tuiApp, ind, structs.BlobNote)
	pageConfirm := "confirmation"
	tuiApp.Pages.AddPage(pageConfirm, delConfirm, true, false)
	tuiApp.Pages.SwitchToPage(pageConfirm)
}
