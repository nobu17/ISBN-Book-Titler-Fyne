package guis

import (
	"errors"
	"fmt"
	"image/color"

	"isbnbook/app/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func getSettingContent() *fyne.Container {
	appSetting := *settings.NewAppSetings()
	appSetting.Init()

	title := widget.NewLabel("Settings")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	caption := widget.NewLabel("Saveボタンを押すと設定が反映されます。\n")
	caption.Wrapping = fyne.TextWrapBreak
	caption.Alignment = fyne.TextAlignLeading

	red := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	errorLabel := canvas.NewText("", red)

	gsPathEntry := makeGSPathEntry(&appSetting)
	zbarPathEntry := makeZBarPathEntry(&appSetting)
	pageSelect := makePageSelect(&appSetting)
	renameSelect := makeRenameSelect(&appSetting)
	readerSelect := makeReaderSelect(&appSetting)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "GS パス", Widget: gsPathEntry},
			{Text: "ZBar パス", Widget: zbarPathEntry},
			{Text: "展開ページ数", Widget: pageSelect},
			{Text: "リネーム設定", Widget: renameSelect},
			{Text: "利用サービス", Widget: readerSelect},
		},
		SubmitText: "Save",
		OnSubmit: func() {
			errorLabel.Text = ""
			err := appSetting.Validate()
			if err != nil {
				errorLabel.Text = fmt.Sprint(err)
			} else {
				appSetting.SaveSetting()
			}
		},
	}
	return container.NewVBox(
		title,
		caption,
		form,
		errorLabel,
	)
}

func makeGSPathEntry(setting *settings.AppSettings) *widget.Entry {
	binding := binding.BindString(&setting.GSPath)
	entry := widget.NewEntryWithData(binding)
	entry.Validator = func(s string) error {
		if s == "" {
			return errors.New("入力してください。")
		}
		return nil
	}
	return entry
}

func makeZBarPathEntry(setting *settings.AppSettings) *widget.Entry {
	binding := binding.BindString(&setting.ZBarPath)
	entry := widget.NewEntryWithData(binding)
	entry.Validator = func(s string) error {
		if s == "" {
			return errors.New("入力してください。")
		}
		return nil
	}
	return entry
}

func makePageSelect(setting *settings.AppSettings) *widget.Select {
	binding := binding.BindString(&setting.ExtractPages)
	selectables := setting.GetSelectablePages()
	sel := widget.NewSelect(selectables, func(s string) {
		binding.Set(s)
	})
	sel.SetSelected(setting.ExtractPages)
	return sel
}

func makeRenameSelect(setting *settings.AppSettings) *widget.Select {
	binding := binding.BindString(&setting.RenameOption)
	selectables := setting.GetSelectableRenames()
	sel := widget.NewSelect(selectables, func(s string) {
		binding.Set(s)
	})
	sel.SetSelected(setting.RenameOption)
	return sel
}

func makeReaderSelect(setting *settings.AppSettings) *fyne.Container {
	readerBind := binding.BindString(&setting.BookReader)
	selectables := setting.GetSelectableReader()

	rakutenEntry := makeRakutenAPIEntry(setting)
	rakutenEntry.Hidden = true

	amazonPAEntry := makeAmazonPAEntry(setting)
	amazonPAEntry.Hidden = true

	var sel *widget.Select
	var allContents *fyne.Container
	sel = widget.NewSelect(selectables, func(s string) {
		readerBind.Set(s)
		if s == settings.RakutenBook.String() {
			rakutenEntry.Hidden = false
		} else {
			rakutenEntry.Hidden = true
		}
		if s == settings.AmazonPA.String() {
			amazonPAEntry.Hidden = false
		} else {
			amazonPAEntry.Hidden = true
		}
		rakutenEntry.Refresh()
		amazonPAEntry.Refresh()
		allContents.Refresh()
	})

	allContents = container.NewVBox(sel)
	allContents.Add(rakutenEntry)
	allContents.Add(amazonPAEntry)

	sel.SetSelected(setting.BookReader)
	return allContents
}

func makeRakutenAPIEntry(setting *settings.AppSettings) *fyne.Container {
	rakutenApikeyBinding := binding.BindString(&setting.RakutenApiKey)
	entry := widget.NewEntryWithData(rakutenApikeyBinding)
	label := widget.NewLabel("API Key:")
	content := container.New(layout.NewBorderLayout(nil, nil, label, nil),
		label, entry)

	return container.NewVBox(content)
}

func makeAmazonPAEntry(setting *settings.AppSettings) *fyne.Container {
	paAssociateIdBinding := binding.BindString(&setting.AmazonPASettings.AssociateId)
	paAccessKeyBinding := binding.BindString(&setting.AmazonPASettings.AccessKey)
	paSecKeyBinding := binding.BindString(&setting.AmazonPASettings.SecretKey)

	identry := widget.NewEntryWithData(paAssociateIdBinding)
	idlabel := widget.NewLabel("AssociateId:")
	idcontent := container.New(layout.NewBorderLayout(nil, nil, idlabel, nil),
		idlabel, identry)
	subContents := container.NewVBox(idcontent)

	paKeyentry := widget.NewEntryWithData(paAccessKeyBinding)
	paKeylabel := widget.NewLabel("AccessKey:  ")
	paKeycontent := container.New(layout.NewBorderLayout(nil, nil, paKeylabel, nil),
		paKeylabel, paKeyentry)
	subContents.Add(paKeycontent)

	secKeyentry := widget.NewEntryWithData(paSecKeyBinding)
	secKeylabel := widget.NewLabel("SecretKey:  ")
	secKeycontent := container.New(layout.NewBorderLayout(nil, nil, secKeylabel, nil),
		secKeylabel, secKeyentry)
	subContents.Add(secKeycontent)

	return subContents
}
