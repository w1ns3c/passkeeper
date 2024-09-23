package tui

import "time"

// RerenderBlobs rerender all windows with blobs
func (tuiApp *TUI) RerenderBlobs() {
	tuiApp.log.Info().
		Msg("start syncing creds")

	defer tuiApp.wg.Done()

	tuiApp.wg.Add(1)
	go func() {
		tuiApp.Usecase.SyncBlobs(tuiApp.Ctx)
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
					Msg("not time to rerender blobs")

				continue
			}

			tuiApp.log.Info().
				Msg("get new list of blobs, time to rerender all")

			// reinit subpages
			tuiApp.SubformCreds = NewCredsList(tuiApp)
			tuiApp.SubformBank = tuiApp.NewBanking(tuiApp.Usecase.GetCards())
			tuiApp.SubformNotes = tuiApp.NewNotes(tuiApp.Usecase.GetNotes())

			tuiApp.SubPages.RemovePage(SubPageCreds)
			tuiApp.SubPages.AddPage(SubPageCreds, tuiApp.SubformCreds, true, false)

			tuiApp.SubPages.RemovePage(SubPageBank)
			tuiApp.SubPages.AddPage(SubPageCreds, tuiApp.SubformCreds, true, false)

			tuiApp.SubPages.RemovePage(SubPageNotes)
			tuiApp.SubPages.AddPage(SubPageCreds, tuiApp.SubformCreds, true, false)

			pageName, _ := tuiApp.SubPages.GetFrontPage()

			tuiApp.SubPages.SwitchToPage(pageName)

		case <-tuiApp.Ctx.Done():
			tuiApp.log.Info().
				Msg("app get signal to down, stopping sync")

			return
		}

	}
}
