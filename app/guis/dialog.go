package guis

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	// "fyne.io/fyne/v2/widget"
)

func NewWaitDialogWithContent(w *fyne.Window, title string, message string) dialog.Dialog {
	// return dialog.NewCustom(title, message, widget.NewProgressBarInfinite(), *w)
	// currently thre is no custom dialog which has no buttons.
	return dialog.NewProgressInfinite(title, message, *w)
}

func NewWaitDialog(w *fyne.Window) dialog.Dialog {
	return NewWaitDialogWithContent(w, "情報", "実行中...")
}
