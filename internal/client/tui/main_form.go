package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FlexMain generate TUI of app main menu (login/register choice)
func FlexMain(tuiApp *TUI) *tview.Flex {

	btn1 := tview.NewButton("Login").
		SetSelectedFunc(func() {
			tuiApp.Pages.SwitchToPage(PageLogin)
		})

	btn1Flex := tview.NewFlex().
		AddItem(nil, 0, 2, false).
		AddItem(btn1, 5, 1, true).
		AddItem(nil, 0, 2, false).
		SetDirection(tview.FlexRow)

	btn2 := tview.NewButton("Register").
		SetSelectedFunc(func() {
			tuiApp.Pages.SwitchToPage(PageRegister)
		})

	btn2Flex := tview.NewFlex().
		AddItem(nil, 0, 2, false).
		AddItem(btn2, 5, 1, true).
		AddItem(nil, 0, 2, false).
		SetDirection(tview.FlexRow)

	flex := tview.NewFlex().
		AddItem(nil, 0, 2, false).
		AddItem(btn1Flex, 10, 1, true).
		AddItem(nil, 0, 1, false).
		AddItem(btn2Flex, 10, 1, true).
		AddItem(nil, 0, 2, false)
	flex.SetTitle("Pass Keeper").SetBorder(true)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		item1 := flex.GetItem(1) // login button
		item2 := flex.GetItem(3) // register button
		switch event.Key() {
		default:
			//tuiApp.Pages.SwitchToPage(PageCreds)
			return event
		case tcell.KeyEsc:
			page := "exit"
			exitModal := ExitModal(tuiApp)
			tuiApp.Pages.AddPage(page, exitModal, true, false)
			tuiApp.Pages.SwitchToPage(page)

		case tcell.KeyTab:

			if !item1.HasFocus() &&
				!btn1Flex.HasFocus() && !btn1.HasFocus() {
				tuiApp.App.SetFocus(btn1)
				return event
			}

			if !item2.HasFocus() &&
				!btn2Flex.HasFocus() && !btn2.HasFocus() {
				tuiApp.App.SetFocus(btn2)
				return event
			}

		case tcell.KeyLeft:

			if !item1.HasFocus() &&
				!btn1Flex.HasFocus() && !btn1.HasFocus() {
				tuiApp.App.SetFocus(btn1)
				return event
			}

			if !item2.HasFocus() &&
				!btn2Flex.HasFocus() && !btn2.HasFocus() {
				tuiApp.App.SetFocus(btn2)
				return event
			}

		case tcell.KeyRight:

			if !item1.HasFocus() &&
				!btn1Flex.HasFocus() && !btn1.HasFocus() {
				tuiApp.App.SetFocus(btn1)
				return event
			}

			if !item2.HasFocus() &&
				!btn2Flex.HasFocus() && !btn2.HasFocus() {
				tuiApp.App.SetFocus(btn2)
				return event
			}

		case tcell.KeyEnter:
			if btn1.HasFocus() {
				tuiApp.Pages.SwitchToPage(PageLogin)
			}
			if btn2.HasFocus() {
				tuiApp.Pages.SwitchToPage(PageRegister)
			}
			return event

		}
		return event
	})

	return flex
}
