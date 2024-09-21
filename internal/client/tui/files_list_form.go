package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"passkeeper/internal/config"
	"passkeeper/internal/entities"
)

var (
	hintFiles = "\n" + genHelp("files")
)

func (tuiApp *TUI) NewFiles(files []*entities.File) *tview.Flex {

	list := NewFileList(files)
	list.Rerender(files)

	list.SetBorder(true).
		SetTitle("Files")

	helpFiles := tview.NewTextView().
		SetTextColor(tcell.ColorBisque).
		SetText(hintFiles)

	flex := tview.NewFlex().
		AddItem(list, 0, 3, true).
		AddItem(helpFiles, 0, 1, false)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// inputs for bank files
		switch event.Key() {
		case tcell.KeyEsc:
			fallthrough
		case tcell.KeyCtrlL:
			pageName := "Logout"
			logoutModal := LogoutModal(tuiApp)
			tuiApp.Pages.AddPage(pageName, logoutModal, true, false)
			tuiApp.Pages.SwitchToPage(pageName)

		case tcell.KeyDelete:
		//list.Delete(tuiApp, ind)
		case tcell.KeyEnter:

			ind := list.GetCurrentItem()
			file, err := tuiApp.Usecase.GetFileByIND(ind)
			if err != nil {
				return nil
			}

			tuiApp.Usecase.StopSync()
			form := tuiApp.NewDownloadForm(file)
			//tuiApp.App.SetFocus(form)
			tuiApp.Pages.AddAndSwitchToPage("download", form, true)

		}

		switch event.Rune() {
		case 'a':
			//tuiApp.Usecase.StopSync()
			//tuiApp.App.SetFocus(viewForm)
			//viewForm.Add(tuiApp, ind, list)
			//NewDownloadForm(tuiApp)

			tuiApp.Usecase.StopSync()
			form := tuiApp.NewUploadForm()
			tuiApp.Pages.AddAndSwitchToPage("upload", form, true)

		case 'e':
			//tuiApp.Usecase.StopSync()
			//tuiApp.App.SetFocus(viewForm)
			//viewForm.Edit(tuiApp, ind, list)

		case 'd':
			//	list.Delete(tuiApp, ind)
		}

		return event
	})

	list.SetChangedFunc(func(ind int, mainText string, secondaryText string, shortcut rune) {
		//card, err := tuiApp.Usecase.GetCardByIND(ind)
		//if err != nil {
		//	tuiApp.log.Error().
		//		Err(err).Msg("wrong card ind")
		//	return
		//}
		//viewForm.Rerender(card)
	})

	return flex
}

type FileList struct {
	*tview.List
	files []*entities.File
}

func NewFileList(files []*entities.File) *FileList {
	list := tview.NewList()
	list.ShowSecondaryText(false).
		SetBorderPadding(0, 0, 0, 0)

	return &FileList{
		List:  list,
		files: files,
	}
}

func (list *FileList) Rerender(files []*entities.File) {
	for ind := list.GetItemCount() - 1; ind >= 0; ind-- {
		list.RemoveItem(ind)
	}

	if files != nil {
		for ind, file := range files {
			res := GenFileShortName(file.Name)
			if ind < 9 {
				list.AddItem(res, "", rune(49+ind), nil)
			} else if ind == 9 {
				list.AddItem(res, "", 'X', nil)
			} else {
				list.AddItem(res, "", rune(65+ind-10), nil)

			}
		}
	}
}

// GenFileShortName beautify file name to show it in the list
func GenFileShortName(filePath string) string {

	var (
		res string
		m   = config.MaxNameLen
	)

	if len(filePath) > m {
		parts := strings.Split(filePath, " ")
		if len(parts) == 1 {
			res = parts[0]
		} else {
			res = strings.Join(parts[:2], "_")
			if len(res) > m {
				res = res[:m]
			}
		}
	} else {
		res = filePath
	}

	return res
}

func (list *FileList) Delete(tuiApp *TUI, ind int) {
	if list.GetItemCount() == 0 {
		return
	}

	delConfirm := DeleteModal(tuiApp, ind, entities.BlobFile)
	pageConfirm := "confirmation"
	tuiApp.Pages.AddPage(pageConfirm, delConfirm, true, false)
	tuiApp.Pages.SwitchToPage(pageConfirm)
}

func (tuiApp *TUI) NewDownloadForm(file *entities.File) tview.Primitive {
	if file == nil {
		return tview.NewForm()
	}

	var (
		btnDownload = "Download"
		btnCancel   = "Cancel"
	)

	modal := tview.NewForm().
		AddInputField("Download folder", "", 0, nil, nil).
		AddButton(btnDownload, nil).
		AddButton(btnCancel, nil)
	modal.SetBorder(true)

	download := func() {
		input := modal.GetFormItem(0).(*tview.InputField)
		dir := input.GetText() // directory to save file inputted in form

		err := tuiApp.Usecase.UnzipAndDownload(dir, file)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("can't unzip/save file blob")

			errModal := NewModalWithParams(tuiApp, "can't unzip/save file blob", PageBlobsMenu)
			tuiApp.Pages.AddAndSwitchToPage("errDownload", errModal, true)

			return
		}

		tuiApp.log.Info().
			Str("id", file.ID).
			Msg("unzip and saved successfully")

		tuiApp.Pages.SwitchToPage(PageBlobsMenu)
		tuiApp.Usecase.ContinueSync()

	}

	cancel := func() {
		tuiApp.Pages.SwitchToPage(PageBlobsMenu)
		tuiApp.Usecase.ContinueSync()
	}

	modal.GetButton(0).SetSelectedFunc(download)
	modal.GetButton(1).SetSelectedFunc(cancel)

	return center(0, 7, modal)
}

func (tuiApp *TUI) NewUploadForm() tview.Primitive {
	var (
		btnUpload = "Upload"
		btnCancel = "Cancel"
	)

	modal := tview.NewForm().
		AddInputField("Uploading file", "", 0, nil, nil).
		AddButton(btnUpload, nil).
		AddButton(btnCancel, nil)
	modal.SetBorder(true)

	upload := func() {
		input := modal.GetFormItem(0).(*tview.InputField)
		filePath := strings.TrimSpace(input.GetText()) // file that need to upload

		file, err := tuiApp.Usecase.ZipAndUpload(filePath)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("can't zip file blob")

			errModal := NewModalWithParams(tuiApp, "can't zip file blob", PageBlobsMenu)
			tuiApp.Pages.AddAndSwitchToPage("errDownload", errModal, true)

			return
		}

		err = tuiApp.Usecase.AddBlob(tuiApp.Ctx, file)
		if err != nil {
			tuiApp.log.Error().
				Err(err).Msg("can't add file blob on server side")

			errModal := NewModalWithParams(tuiApp, "can't upload file blob", PageBlobsMenu)
			tuiApp.Pages.AddAndSwitchToPage("errDownload", errModal, true)

			return
		}

		tuiApp.log.Error().
			Str("id", file.ID).
			Msg("add file blob successfully")

		tuiApp.Pages.SwitchToPage(PageBlobsMenu)
		tuiApp.Usecase.ContinueSync()

	}

	cancel := func() {
		tuiApp.Pages.SwitchToPage(PageBlobsMenu)
		tuiApp.Usecase.ContinueSync()
	}

	modal.GetButton(0).SetSelectedFunc(upload)
	modal.GetButton(1).SetSelectedFunc(cancel)

	return center(0, 7, modal)
}

// center returns a new primitive which shows the provided primitive in its
// center, given the provided primitive's size.
func center(width, height int, p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
