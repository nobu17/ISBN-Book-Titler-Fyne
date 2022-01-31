package workflows

import (
	"fmt"
	"os"

	"isbnbook/app/log"
	"isbnbook/app/settings"
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
	if preResult := w.preCheck(path); preResult != nil {
		return preResult
	}

	isbnflow := isbn.NewIsbnGetWorkFlow(w.appSettings)
	isbn13, err := isbnflow.GetIsbn(path, int(w.appSettings.GetPagesInt()))
	if err != nil {
		logger.Error("Get ISBN Error. err:", err)
		return &WorkFlowResult{"", path, "Get ISBN Error:", err}
	}

	info, err := book.GetBookInfo(isbn13, w.appSettings)
	if err != nil {
		logger.Error("Get Bookinfo Error. err:", err)
		return &WorkFlowResult{"", path, "Get Bookinfo Error:", err}
	}

	newname, err := rename.Rename(path, w.appSettings, w.ruleSetings, info)
	if err != nil {
		logger.Error("Rename file Error. err:", err)
		return &WorkFlowResult{"", path, "Rename file Error:", err}
	}
	return &WorkFlowResult{newname, path, "", nil}
}

func (w *RenameByBookInfoWorkflow) TestGetBookInfo(path string) (*book.BookInfo, error) {
	if preResult := w.preCheck(path); preResult != nil {
		return nil, preResult.Error
	}
	isbnflow := isbn.NewIsbnGetWorkFlow(w.appSettings)
	isbn13, err := isbnflow.GetIsbn(path, int(w.appSettings.GetPagesInt()))
	if err != nil {
		return nil, err
	}

	return book.GetBookInfo(isbn13, w.appSettings)
}

func (w *RenameByBookInfoWorkflow) preCheck(path string) *WorkFlowResult {
	if err := w.appSettings.Validate(); err != nil {
		errmsg := "appsetting is incorect"
		logger.Error(errmsg, err)
		return &WorkFlowResult{"", path, errmsg, err}
	}
	if err := w.ruleSetings.Validate(); err != nil {
		errmsg := "ruleSetings is incorect"
		logger.Error(errmsg, err)
		return &WorkFlowResult{"", path, errmsg, err}
	}

	if !w.isFileExists(path) {
		logger.Error("file is not exists", nil)
		return &WorkFlowResult{"", path, "file is not exists", fmt.Errorf("file is not exists")}
	}
	return nil
}

func (w *RenameByBookInfoWorkflow) isFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
