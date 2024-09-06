package tui

import "time"

func (tuiApp *TUI) RerenderCreds() {
	go tuiApp.Usecase.SyncCreds(tuiApp.Ctx)

	t := tuiApp.Usecase.GetSyncTime()
	ticker := time.NewTicker(t)
	for range ticker.C {
		credsForm := NewCredsList(tuiApp)
		tuiApp.Pages.RemovePage(PageCreds)
		tuiApp.Pages.AddPage(PageCreds, credsForm, true, false)
		tuiApp.Pages.SwitchToPage(PageCreds)
	}

}
