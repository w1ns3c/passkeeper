package tui

import (
	"time"

	"github.com/rivo/tview"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/structs"
)

type NoteDetails struct {
	*tview.Form

	//tuiApp *TUI

	FieldName *tview.InputField
	FieldDate *tview.InputField
	FieldBody *tview.TextArea

	CurrentNote *structs.Note

	// fields sizes
	FieldWidth  int
	FieldHeight int
	maxSigns    int
}

func NewNoteDetails(note *structs.Note) *NoteDetails {
	if note == nil {
		note = &structs.Note{}
	}

	form := &NoteDetails{
		Form:        tview.NewForm(),
		FieldName:   tview.NewInputField().SetLabel("Name:").SetText(note.Name),
		FieldDate:   tview.NewInputField().SetLabel("Date").SetText(note.Date.Format(time.DateTime)),
		FieldBody:   tview.NewTextArea().SetLabel("Note:").SetText(note.Body, true),
		CurrentNote: note,
		FieldWidth:  40,
		FieldHeight: 6,
		maxSigns:    config.MaxNoteLen,
	}

	form.Form.SetBorder(true).
		SetTitle("Details")

	form.FieldDate.SetFieldWidth(form.FieldWidth).
		SetDisabled(true)

	form.FieldBody.SetMaxLength(form.maxSigns)

	form.Form.AddFormItem(form.FieldName)
	form.Form.AddFormItem(form.FieldDate)
	form.Form.AddFormItem(form.FieldBody)

	return form
}

func (form *NoteDetails) Rerender(note *structs.Note) {
	form.FieldName.SetText(note.Name)
	form.FieldDate.SetText(note.Date.Format(time.DateTime))
	form.FieldBody.SetText(note.Body, true)
	form.CurrentNote = note
}

// Add responsible for TUI of adding new entities.Note
func (form *NoteDetails) Add(tuiApp *TUI, ind int, list *NotesList) {
	form.EmptyFields()
	form.AddButton("Save", func() {
		newNote, err := form.GetCurrentValues()
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add new note on client side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageBlobUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageBlobUpdError)
			return
		}

		newNote.ID = hashes.GeneratePassID()

		if err := tuiApp.Usecase.AddBlob(tuiApp.Ctx, newNote); err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add new note on server side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageBlobUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageBlobUpdError)
			return
		}

		// rerender credsList
		list.Rerender(tuiApp.Usecase.GetNotes())
		form.Rerender(newNote)
		tuiApp.App.SetFocus(list)

		// hide buttons
		form.HideButtons()

		//defer continue cred sync
		tuiApp.Usecase.ContinueSync()

	})

	form.AddButton("Cancel", func() {
		//defer continue cred sync
		defer tuiApp.Usecase.ContinueSync()

		// defer remove buttons
		defer form.HideButtons()
		// rerender credsList
		defer tuiApp.App.SetFocus(list)
		if tuiApp.Usecase.CredsLen() > 0 {
			note, _ := tuiApp.Usecase.GetNoteByIND(ind)
			form.Rerender(note)

			return
		}

		// clear fields if there isn't any blobsUC
		form.HideFields()

		curNote, err := tuiApp.Usecase.GetNoteByIND(list.GetCurrentItem())
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("can't get current note")

			return
		}

		form.Rerender(curNote)
	})
}

// Edit responsible for TUI of editing current selected entities.Note
func (form *NoteDetails) Edit(tuiApp *TUI, ind int, list *NotesList) {
	tmp := tuiApp.Usecase.GetNotes()
	if tmp == nil || len(tmp) <= ind || len(tmp) == 0 {
		return
	}

	form.AddButton("Save", func() {
		cur, err := tuiApp.Usecase.GetNoteByIND(ind)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to edit note on client side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageBlobUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageBlobUpdError)

			return
		}

		editedNote, err := form.GetCurrentValues()
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to edit note on client side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageBlobUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageBlobUpdError)

			return
		}

		editedNote.SetID(cur.GetID())

		if err := tuiApp.Usecase.EditBlob(tuiApp.Ctx, editedNote, ind); err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to edit note on server side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageBlobUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageBlobUpdError)

			return
		}

		// rerender credsList
		list.Rerender(tuiApp.Usecase.GetNotes())
		form.Rerender(editedNote)
		tuiApp.App.SetFocus(list)

		// hide buttons
		form.HideButtons()

		//defer continue cred sync
		tuiApp.Usecase.ContinueSync()

	})

	form.AddButton("Cancel", func() {
		//defer continue cred sync
		defer tuiApp.Usecase.ContinueSync()

		// defer remove buttons
		defer form.HideButtons()
		// rerender credsList
		defer tuiApp.App.SetFocus(list)
		if tuiApp.Usecase.CredsLen() > 0 {
			note, err := tuiApp.Usecase.GetNoteByIND(ind)
			if err != nil {
				tuiApp.log.Error().
					Err(err).Msg("can't get current note")

				return
			}
			form.Rerender(note)

			return
		}

		// clear fields if there isn't any blobsUC
		//form.EmptyFields()
		form.HideFields()
	})
}

// EmptyFields reset fields' values
func (form *NoteDetails) EmptyFields() {
	if form.FieldName != nil {
		form.FieldName.SetText("")
	}

	if form.FieldDate != nil {
		form.FieldDate.SetText("")
	}

	if form.FieldBody != nil {
		form.FieldBody.SetText("", true)
	}

}

// GetCurrentValues get values from user input and format it to Note entity
func (form *NoteDetails) GetCurrentValues() (newNote *structs.Note, err error) {
	newNote = new(structs.Note)
	newNote.Type = structs.BlobNote

	newNote.Name = form.FieldName.GetText()
	newNote.Date = time.Now()
	newNote.Body = form.FieldBody.GetText()

	return newNote, nil
}

// HideButtons hide Save/Cancel buttons
func (form *NoteDetails) HideButtons() {
	if form.GetButtonCount() <= 0 {
		return
	}

	form.RemoveButton(1)
	form.RemoveButton(0)
}

// HideFields remove all items from form
func (form *NoteDetails) HideFields() {
	for ind := form.Form.GetFormItemCount() - 1; ind >= 0; ind-- {
		form.Form.RemoveFormItem(ind)
	}
}
