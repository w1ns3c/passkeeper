package tui

import (
	"github.com/rivo/tview"

	"passkeeper/internal/entities"
)

type NoteDetails struct {
	*tview.Form

	//tuiApp *TUI

	FieldName *tview.InputField
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
		Form: tview.NewForm(),
		//tuiApp:          nil,
		FieldName:   tview.NewInputField().SetLabel("Name:").SetText(note.Name),
		FieldBody:   tview.NewTextArea().SetLabel("Note:").SetText(note.Body, true),
		BtnSave:     nil,
		BtnCancel:   nil,
		SaveLabel:   "",
		CancelLabel: "",
		CurrentNote: note,
		FieldWidth:  0,
		FieldHeight: 0,
		maxSigns:    0,
	}

	form.Form.SetBorder(true).
		SetTitle("Details")

	form.Form.AddFormItem(form.FieldName)
	form.Form.AddFormItem(form.FieldBody)

	return form
}

func (form *NoteDetails) Rerender(note *entities.Note) {
	form.FieldName.SetText(note.Name)
	form.FieldBody.SetText(note.Body, true)
}
