package guis

import (
	"fmt"

	"isbnbook/app/settings"
	"isbnbook/app/workflows"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var data = [][]string{
	{"結果", "旧ファイル名", "新ファイル名", "エラー詳細"},
}

type taskCallback func([]workflows.WorkFlowResult)

func getResultContent(w *fyne.Window, files []string) *fyne.Container {

	title := widget.NewLabel("Results")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	caption := widget.NewLabel("現在のバージョンは画面から実行できません。\n一回アプリを閉じて、実行ファイル(exe)にディレクトリかファイルをドラッグ&ドロップしてください。\n")
	caption.Wrapping = fyne.TextWrapBreak
	caption.Alignment = fyne.TextAlignLeading

	top := container.NewVBox(
		title,
		caption,
	)
	table := makeTables()
	content := container.New(layout.NewBorderLayout(top, nil, nil, nil),
		top, table)

	if len(files) > 0 {
		diag := NewWaitDialog(w)
		diag.Show()
		go executeWorkFlow(files, diag, func(result []workflows.WorkFlowResult) {
			updateTable(result, table)
			diag.Hide()
		})
	}

	return content
}

func executeWorkFlow(files []string, diag dialog.Dialog, cb taskCallback) {
	appsetting := settings.NewAppSetings()
	appsetting.Init()

	rulesetting := settings.NewRuleSettings()
	rulesetting.Init()

	results := []workflows.WorkFlowResult{}

	flow := workflows.NewRenameByBookInfoWorkflow(appsetting, rulesetting)
	for _, file := range files {
		result := flow.RenameFileByIsbn(file)
		results = append(results, *result)
	}
	cb(results)
}

func updateTable(results []workflows.WorkFlowResult, table *widget.Table) {
	for _, item := range results {
		col1 := "成功"
		col4 := "なし"
		if item.Error != nil {
			col1 = "失敗"
			col4 = fmt.Sprintf("%s %s", item.Message, item.Error.Error())
		}
		col2 := item.OldName
		col3 := item.NewName
		row := []string{col1, col2, col3, col4}
		data = append(data, row)
	}
	table.Refresh()
}

func makeTables() *widget.Table {
	var t *widget.Table
	colwidth := make([]float32, len(data[0]))
	t = widget.NewTable(
		func() (int, int) { return len(data), len(data[0]) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			label.SetText(fmt.Sprintf(data[id.Row][id.Col]))
			// adjust column size
			currentWidth := label.MinSize().Width
			if colwidth[id.Col] < currentWidth {
				colwidth[id.Col] = currentWidth
				go func() {
					t.SetColumnWidth(id.Col, colwidth[id.Col])
				}()
			}
		})
	return t
}
