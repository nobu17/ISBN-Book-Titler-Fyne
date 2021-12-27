package guis

import (
	"errors"
	"fmt"
	"image/color"
	
	"isbnbook/app/settings"
	"isbnbook/app/workflows/rename"
	"isbnbook/app/log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func GetRuleContent() *fyne.Container {
	logger := log.GetLogger()

	rule := *settings.NewRuleSettings()
	rule.Init()

	title := widget.NewLabel("Rules")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	caption := widget.NewLabel("Saveボタンを押すと設定が反映されます。\n")
	caption.Wrapping = fyne.TextWrapBreak
	caption.Alignment = fyne.TextAlignLeading

	red := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	errorLabel := canvas.NewText("", red)

	renameEntry := makeReNameEntry(&rule)
	renameEntry.Validator = func(s string) error {
		if s == "" {
			return errors.New("入力してください。")
		}

		return nil
	}

	explanLabel := widget.NewLabel(rename.GetExplaination())
	explanLabel.Wrapping = fyne.TextWrapBreak
	explanLabel.TextStyle.Bold = true

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ReNnme Rules", Widget: renameEntry},
		},
		SubmitText: "Save",
		OnSubmit: func() {
			errorLabel.Text = ""
			if err := rule.Validate(); err != nil {
				errorLabel.Text = fmt.Sprint(err)
				return
			}
			if err := rule.SaveSetting(); err != nil {
				errorLabel.Text = "save failed...." + fmt.Sprint(err)
				logger.Error("save failed.", err)
			}
		},
	}
	return container.NewVBox(
		title,
		caption,
		form,
		errorLabel,
		explanLabel,
	)
}

func makeReNameEntry(setting *settings.RuleSettings) *widget.Entry {
	binding := binding.BindString(&setting.RenameRule)
	return widget.NewEntryWithData(binding)
}
