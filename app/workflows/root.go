package workflows

import (
	"fmt"
	"os"

	"isbnbook/app/settings"
	"isbnbook/app/log"
	"isbnbook/app/workflows/book"
	"isbnbook/app/workflows/isbn"
	"isbnbook/app/workflows/rename"
)

type RenameByBookInfoWorkflow struct {
	appSettings *settings.AppSettings
	ruleSetings *settings.RuleSettings
}

type WorkFlowResult struct {
	NewName string
	OldName string
	Message string
	Error   error
}

var logger = log.GetLogger()

func NewRenameByBookInfoWorkflow(app *settings.AppSettings, rule *settings.RuleSettings) *RenameByBookInfoWorkflow {
	return &RenameByBookInfoWorkflow{
		app,
		rule,
	}
}

func (w *RenameByBookInfoWorkflow) RenameFileByIsbn(path string) *WorkFlowResult {
	if !w.isFileExists(path) {
		logger.Error("file is not exists", nil)
		return &WorkFlowResult{"", path, "file is not exists", fmt.Errorf("file is not exists")}
	}
	isbnflow := isbn.NewIsbnGetWorkFlow(w.appSettings)
	isbn13, err := isbnflow.GetIsbn(path, int(w.appSettings.GetPagesInt()))
	if err != nil {
		logger.Error("Get ISBN Error", err)
		return &WorkFlowResult{"", path, "Get ISBN Error", err}
	}

	info, err := book.GetBookInfo(isbn13, w.appSettings)
	if err != nil {
		logger.Error("Get Bookinfo Error", err)
		return &WorkFlowResult{"", path, "Get Bookinfo Error", err}
	}

	newname, err := rename.Rename(path, w.ruleSetings, info)
	if err != nil {
		logger.Error("Rename file Error", err)
		return &WorkFlowResult{"", path, "Rename file Error", err}
	}
	return &WorkFlowResult{newname, path, "", nil}
}

func (w *RenameByBookInfoWorkflow) TestGetBookInfo(path string) (*book.BookInfo, error) {
	if !w.isFileExists(path) {
		return nil, fmt.Errorf("file is not exists")
	}
	isbnflow := isbn.NewIsbnGetWorkFlow(w.appSettings)
	isbn13, err := isbnflow.GetIsbn(path, int(w.appSettings.GetPagesInt()))
	if err != nil {
		return nil, err
	}

	return book.GetBookInfo(isbn13, w.appSettings)
}

func (w *RenameByBookInfoWorkflow) isFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
