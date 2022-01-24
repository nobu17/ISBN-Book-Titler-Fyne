package rename

import (
	"isbnbook/app/log"
	"isbnbook/app/settings"
	"isbnbook/app/workflows/book"
)

var manager = newRenameManager()
var logger = log.GetLogger()

func Rename(filePath string, appSetting *settings.AppSettings, renameSetting *settings.RuleSettings, bookInfo *book.BookInfo) (string, error) {
	repname, err := manager.GetReplaceName(renameSetting, bookInfo)
	if err != nil {
		logger.Error("get rename error", err)
		return "", err
	}
	newname, err := manager.Rename(filePath, repname, appSetting)
	if err != nil {
		logger.Error("rename error", err)
		return "", err
	}
	return newname, err
}

func GetExplaination() string {
	return manager.GetExplaination()
}
