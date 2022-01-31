package main

import (
	"isbnbook/app/command"
	"isbnbook/app/guis"
	"isbnbook/app/log"
	mytheme "isbnbook/app/theme"
	"isbnbook/app/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type menu struct {
	Name        string
	Icon        fyne.Resource
	GetMenuItem func() *fyne.Container
}

func createMenus(w *fyne.Window, files []string) *[]menu {
	menuItems := make([]menu, 0)

	// setting
	settings := guis.NewSettingDisplay()
	settingMenu := menu{"Settings", theme.SettingsIcon(), func() *fyne.Container {
		return settings.Display(w, nil)
	}}
	menuItems = append(menuItems, settingMenu)

	// rule
	rule := guis.NewRuleDisplay()
	ruleMenu := menu{"Rules", theme.DocumentIcon(), func() *fyne.Container {
		return rule.Display(w, nil)
	}}
	menuItems = append(menuItems, ruleMenu)

	// test
	test := guis.NewTestDisplay()
	testMenu := menu{"Test", theme.CancelIcon(), func() *fyne.Container {
		return test.Display(w, nil)
	}}
	menuItems = append(menuItems, testMenu)

	// results
	result := guis.NewResultsDisplay()
	resultMenu := menu{"Results", theme.ComputerIcon(), func() *fyne.Container {
		return result.Display(w, files)
	}}
	menuItems = append(menuItems, resultMenu)

	// version
	version := guis.NewVersionDisplay()
	versionMenu := menu{"Version", theme.DownloadIcon(), func() *fyne.Container {
		return version.Display(w, nil)
	}}
	menuItems = append(menuItems, versionMenu)

	return &menuItems
}

func main() {
	logger := log.GetLogger()
	logger.Info("start app")

	// change work dir
	err := utils.ChangeWorkDir()
	if err != nil {
		logger.Error("failed to init change work dir", err)
		return
	}

	// get file or dir from drag and drop
	files, failed, err := command.GetFiles()
	if err != nil {
		logger.Error("error from get files", err)
		return
	}

	a := app.New()
	a.Settings().SetTheme(&mytheme.MyTheme{})
	w := a.NewWindow("ISBN Book Titler Fyne")
	// a.Settings().SetTheme(theme.LightTheme())

	// creating a content
	body := container.NewVBox()

	// init meenu
	menuItems := *createMenus(&w, files)

	list := widget.NewList(
		func() int {
			return len(menuItems)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(menuItems[id].Name)
			item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(menuItems[id].Icon)
		},
	)
	// when selecting a left side menu, displaying right side contents
	list.OnSelected = func(id widget.ListItemID) {
		left := list
		body = menuItems[id].GetMenuItem()
		content := container.New(
			layout.NewBorderLayout(nil, nil, left, nil),
			left, body)

		w.SetContent(content)
	}

	if len(files) > 0 {
		list.Select(3)
	} else {
		list.Select(0)
	}
	// check command result
	if len(failed) > 0 && len(files) == 0 {
		dialog.ShowInformation("error", "有効なファイルがありませんでした。", w)
	}

	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
