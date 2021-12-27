package guis

import (
	"image/color"
	"strings"

	"isbnbook/app/settings"
	"isbnbook/app/workflows"
	"isbnbook/app/workflows/book"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const TestFilePath = "./samples/test.pdf"

func GetTestContent(w *fyne.Window) *fyne.Container {

	title := widget.NewLabel("Test")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	// caption := widget.NewLabel("Settingsタブの内容でテストを行います。\nGSやZbarのパスが正しく設定されているか確認できます。\nStart Testボタンを押すと下記のパスにあるファイルから画像を取り出して書籍情報を取得します。\n1:ファイルを差し替えれば自分のファイルでも確認可能です。\n2:Saveボタンを押していない設定は使用されません。\n" + TestFilePath)
	// caption := widget.NewLabel("aaa")
	captions := createCaption()
	// caption.Wrapping = fyne.TextWrapWord
	//caption.Alignment = fyne.TextAlignLeading

	bookmsg := ""
	bookbind := binding.BindString(&bookmsg)
	bookLabel := widget.NewLabelWithData(bookbind)
	bookLabel.Wrapping = fyne.TextWrapBreak
	bookLabel.TextStyle.Bold = true

	red := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	errorLabel := canvas.NewText("", red)

	testButton := widget.NewButton("Start Test", func() {
		rulesetting := settings.NewRuleSettings()
		rulesetting.Init()

		appsetting := settings.NewAppSetings()
		appsetting.Init()

		diag := NewWaitDialog(w)
		diag.Show()
		defer diag.Hide()

		bookbind.Set("")
		flow := workflows.NewRenameByBookInfoWorkflow(appsetting, rulesetting)
		book, err := flow.TestGetBookInfo(TestFilePath)
		if err != nil {
			errorLabel.Text = "failed to read book info." + err.Error()
			return
		}
		bookbind.Set(getBookDisplayInfo(book))
	})

	return container.NewVBox(
		title,
		captions,
		testButton,
		errorLabel,
		bookLabel,
	)
}

func createCaption() *fyne.Container {
	captions := container.NewVBox(
		widget.NewLabel("Settingsタブの内容でテストを行います。(Saveボタンを押していない設定は使用されません。)"),
		widget.NewLabel("GSやZbarのパスが正しく設定されているか確認できます。"),
		widget.NewLabel("ファイルを差し替えれば自分のファイルでも確認可能です。"+TestFilePath),
	)
	return captions
}

func getBookDisplayInfo(bookinfo *book.BookInfo) string {
	str := ""
	str += "title:" + bookinfo.Title + "\n"
	str += "authors:" + strings.Join(bookinfo.Authors, ",") + "\n"
	str += "publisher:" + bookinfo.Publisher + "\n"
	str += "date:" + bookinfo.Date + "\n"
	str += "genre:" + bookinfo.Genre + "\n"
	str += "kind:" + bookinfo.Kind
	return str
}
