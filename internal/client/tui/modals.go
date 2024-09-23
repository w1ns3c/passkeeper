package tui

import (
	"github.com/rivo/tview"

	"passkeeper/internal/entities/structs"
)

func DeleteModal(tuiApp *TUI, ind int, blobType structs.BlobType) *tview.Modal {
	btn1Name := "Yes"
	btn2Name := "No"

	errModal := tview.NewModal().
		SetText("Are you really want to delete this item?").
		AddButtons([]string{btn1Name, btn2Name}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == btn1Name {
				err := tuiApp.Usecase.DelBlob(tuiApp.Ctx, ind, blobType)
				if err != nil {
					tuiApp.log.Error().
						Err(err).Msg("can't delete cred on server side")

					errModal := NewErrorEditModal(tuiApp, "can't delete cred on server side", tuiApp.Pages)
					tuiApp.Pages.AddAndSwitchToPage("pageDel", errModal, true)

					return
				}

				credsForm := NewCredsList(tuiApp)
				tuiApp.Pages.RemovePage(SubPageCreds)
				tuiApp.Pages.AddPage(SubPageCreds, credsForm, true, false)

				pageDel := "deleted"
				deletedPage := NewModalWithParams(tuiApp, "Credential successful deleted!", PageBlobsMenu)
				tuiApp.Pages.AddPage(pageDel, deletedPage, true, false)
				tuiApp.Pages.SwitchToPage(pageDel)

				tuiApp.Rerender()
			}

			if buttonLabel == btn2Name {
				tuiApp.Pages.SwitchToPage(PageBlobsMenu)
			}
		}).
		SetFocus(1)

	return errModal
}

func ExitModal(tuiApp *TUI) *tview.Modal {
	btn1Name := "Yes"
	btn2Name := "No"

	exitModal := tview.NewModal().
		SetText("Are you really want to exit?").
		AddButtons([]string{btn1Name, btn2Name}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == btn1Name {
				err := tuiApp.interruptSignal()
				if err != nil {
					tuiApp.log.Error().
						Err(err).Msg("cannot send interrupt signal")

					return
				}
			}
			if buttonLabel == btn2Name {
				tuiApp.Pages.SwitchToPage(PageMain)
			}
		}).
		SetFocus(1)

	return exitModal

}

func LogoutModal(tuiApp *TUI) *tview.Modal {
	btn1Name := "Logout"
	btn2Name := "Cancel"

	logoutModal := tview.NewModal().
		SetText("Are you really want to logout?").
		AddButtons([]string{btn1Name, btn2Name}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == btn1Name {
				tuiApp.Logout()
				tuiApp.Pages.SwitchToPage(PageMain)
				tuiApp.App.SetFocus(tuiApp.FormAuth)
			}
			if buttonLabel == btn2Name {
				tuiApp.Pages.SwitchToPage(PageBlobsMenu)
			}
		}).
		SetFocus(1)

	return logoutModal
}

func NewModalWithParams(tuiApp *TUI, text string, page string) *tview.Modal {
	text = CapitalizeFirst(text)
	errModal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				tuiApp.Pages.SwitchToPage(page)
				return
			}
		})

	return errModal
}

func NewErrorEditModal(tuiApp *TUI, text string, returnPage tview.Primitive) *tview.Modal {
	text = CapitalizeFirst(text)
	errModal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				tuiApp.Pages.SwitchToPage(PageBlobsMenu)
				tuiApp.App.SetFocus(returnPage)
				return
			}
		})

	return errModal
}

// CapitalizeFirst capitalize the first letter in text
func CapitalizeFirst(text string) string {
	// Capitalize first letter
	if len(text) != 0 {
		sign := text[0]
		if 'a' < sign && sign < 'z' {
			sign -= 32
		}
		text = string(sign) + text[1:]

	}
	return text
}
