package tui

import (
	"time"

	"github.com/rivo/tview"

	"passkeeper/internal/entities"
	"passkeeper/internal/entities/hashes"
)

type NoteDetails struct {
	*tview.Form

	//tuiApp *TUI

	FieldName *tview.InputField
	FieldDate *tview.InputField
	FieldBody *tview.TextArea

	BtnSave   *tview.Button
	BtnCancel *tview.Button

	SaveLabel   string
	CancelLabel string

	CurrentNote *entities.Note

	// fields sizes
	FieldWidth  int
	FieldHeight int
	maxSigns    int
}

func NewNoteDetails(note *entities.Note) *NoteDetails {
	if note == nil {
		note = &entities.Note{}
	}

	form := &NoteDetails{
		Form:      tview.NewForm(),
		FieldName: tview.NewInputField().SetLabel("Name:").SetText(note.Name),
		FieldDate: tview.NewInputField().SetLabel("Date").
			SetText(note.Date.Format(time.DateTime)),
		FieldBody:   tview.NewTextArea().SetLabel("Note:").SetText(note.Body, true),
		BtnSave:     nil,
		BtnCancel:   nil,
		SaveLabel:   "",
		CancelLabel: "",
		CurrentNote: note,
		FieldWidth:  40,
		FieldHeight: 6,
		maxSigns:    0,
	}

	form.Form.SetBorder(true).
		SetTitle("Details")

	form.FieldDate.SetFieldWidth(form.FieldWidth).
		SetDisabled(true)

	form.Form.AddFormItem(form.FieldName)
	form.Form.AddFormItem(form.FieldDate)
	form.Form.AddFormItem(form.FieldBody)

	return form
}

func (form *NoteDetails) Rerender(note *entities.Note) {
	form.FieldName.SetText(note.Name)
	form.FieldBody.SetText(note.Body, true)
	form.CurrentNote = note
}

func (form *NoteDetails) Add(tuiApp *TUI, ind int, list *NotesList) {
	form.EmptyFields()
	form.AddButton("Save", func() {
		newNote, err := form.GetCurrentValues()
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add new note on client side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)
			return
		}

		newNote.ID = hashes.GeneratePassID2()

		if err := tuiApp.Usecase.AddBlob(tuiApp.Ctx, newNote); err != nil {
			tuiApp.log.Error().
				Err(err).Msg("failed to add new note on server side")
			errModal := NewErrorEditModal(tuiApp, err.Error(), form)
			tuiApp.Pages.AddPage(PageCredUpdError, errModal, true, false)
			tuiApp.Pages.SwitchToPage(PageCredUpdError)
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
func (form *NoteDetails) GetCurrentValues() (newNote *entities.Note, err error) {
	newNote = new(entities.Note)
	newNote.Type = entities.UserNote

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
