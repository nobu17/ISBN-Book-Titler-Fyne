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

func GetSettingContent() *fyne.Container {
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
	readerSelect := makeReaderSelect(&appSetting)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "GS パス", Widget: gsPathEntry},
			{Text: "ZBar パス", Widget: zbarPathEntry},
			{Text: "展開ページ数", Widget: pageSelect},
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

func makeReaderSelect(setting *settings.AppSettings) *fyne.Container {
	readerBind := binding.BindString(&setting.BookReader)
	selectables := setting.GetSelectableReader()
	rakutenApikeyBinding := binding.BindString(&setting.RakutenApiKey)

	var sel *widget.Select
	var allContents *fyne.Container
	sel = widget.NewSelect(selectables, func(s string) {
		readerBind.Set(s)
		if s == settings.RakutenBook.String() {
			entry := widget.NewEntryWithData(rakutenApikeyBinding)
			label := widget.NewLabel("API Key:")
			content := container.New(layout.NewBorderLayout(nil, nil, label, nil),
			label, entry)
			allContents.Add(content)
		}
	})

	allContents = container.NewVBox(sel)


	sel.SetSelected(setting.BookReader)
	return allContents
}
