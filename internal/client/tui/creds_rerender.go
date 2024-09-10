package tui

import "time"

func (tuiApp *TUI) RerenderCreds() {
	go tuiApp.Usecase.SyncCreds(tuiApp.Ctx)

	t := tuiApp.Usecase.GetSyncTime()
	ticker := time.NewTicker(t)
	for {
		select {
		case <-ticker.C:
			// don't rerender forms if user edit/add new cred in tui
			if tuiApp.Usecase.CheckSync() {
				continue
			}

			credsForm := NewCredsList(tuiApp)
			tuiApp.Pages.RemovePage(PageCreds)
			tuiApp.Pages.AddPage(PageCreds, credsForm, true, false)
			tuiApp.Pages.SwitchToPage(PageCreds)
		case <-tuiApp.Ctx.Done():
			return
		}

	}
}
