package guis

import (
	"fmt"
	"net/url"

	"isbnbook/app/versions"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func getVersionContent(w *fyne.Window) *fyne.Container {

	title := widget.NewLabel("Version")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	current := ""
	currentBind := binding.BindString(&current)
	currentLabel := widget.NewLabelWithData(currentBind)
	currentLabel.Alignment = fyne.TextAlignCenter

	latest := ""
	latestBind := binding.BindString(&latest)
	latestLabel := widget.NewLabelWithData(latestBind)
	latestLabel.Alignment = fyne.TextAlignCenter

	url ,_ := url.Parse("https://github.com/nobu17/ISBN-Book-Titler-Fyne")
	link := widget.NewHyperlink("最新版をDLする", url)

	testFunc := func() {

		diag := NewWaitDialog(w)
		diag.Show()
		defer diag.Hide()

		ver, err := versions.NewClient()
		if err != nil {
			currentBind.Set(getCurrentString("エラーが発生しました。"))
			latestBind.Set(getLatestString("エラーが発生しました。"))
			return
		}

		currentBind.Set(getCurrentString(ver.GetCurrent().String()))

		latestVer, err := ver.GetLatest()
		if err != nil {
			latestBind.Set(getLatestString("エラーが発生しました。"))
			return
		}
		latestBind.Set(getLatestString(latestVer.String()))
	}

	testFunc()

	reloadBtn := widget.NewButton("Reload", func() {
		testFunc()
	})

	return container.NewVBox(
		title,
		currentLabel,
		latestLabel,
		reloadBtn,
		link,
	)
}

func getCurrentString(content string) string {
	return fmt.Sprintf("使用 Version: %s", content)
}

func getLatestString(content string) string {
	return fmt.Sprintf("最新 Version: %s", content)
}
