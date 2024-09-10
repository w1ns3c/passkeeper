package tui

import "time"

func (tuiApp *TUI) RerenderCreds() {
	tuiApp.log.Info().
		Msg("start syncing creds")

	defer tuiApp.wg.Done()

	tuiApp.wg.Add(1)
	go func() {
		tuiApp.Usecase.SyncCreds(tuiApp.Ctx)
		tuiApp.wg.Done()
	}()

	t := tuiApp.Usecase.GetSyncTime()
	ticker := time.NewTicker(t)
	for {
		select {
		case <-ticker.C:
			// don't rerender forms if user edit/add new cred in tui
			if tuiApp.Usecase.CheckSync() {
				tuiApp.log.Info().
					Msg("not time to sync creds: user changing them")

				continue
			}

			tuiApp.log.Info().
				Msg("get new list of creds, time to rerender all")

			credsForm := NewCredsList(tuiApp)
			tuiApp.Pages.RemovePage(PageCreds)
			tuiApp.Pages.AddPage(PageCreds, credsForm, true, false)
			tuiApp.Pages.SwitchToPage(PageCreds)
		case <-tuiApp.Ctx.Done():
			tuiApp.log.Info().
				Msg("app get signal to down, stopping sync")

			return
		}

	}
}
