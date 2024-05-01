package forms

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateForm() {

	form := tview.NewForm().
		AddInputField("tmp", "123", 0, nil, nil).
		AddButton("Login", tmp).
		AddButton("Register", tmp).SetTitle("Auth").
		SetBorder(true)

	app := tview.NewApplication().SetRoot(form, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyCtrlC:
			return nil
		case tcell.KeyDown:
			//app.SetFocus()
		}
		return event
	})
	app.Run()

}