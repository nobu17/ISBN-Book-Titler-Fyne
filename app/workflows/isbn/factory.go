package isbn

import (
	"fmt"
	"strings"

	"isbnbook/app/settings"
)

type fileExtractor interface {
	Extract(workDir string) error
}

func NewFileExtractor(settings *settings.AppSettings, filePath string) (fileExtractor, error) {
	if strings.HasSuffix(strings.ToLower(filePath), ".pdf") {
		return newPDFExtractor(settings.GSPath, filePath, int(settings.GetPagesInt())), nil
	} else if strings.HasSuffix(strings.ToLower(filePath), ".zip") {
		return newZipExtractor(filePath, int(settings.GetPagesInt())), nil
	}  else if strings.HasSuffix(strings.ToLower(filePath), ".rar") {
		return newRarExtractor(filePath, int(settings.GetPagesInt())), nil
	} else {
		return nil, fmt.Errorf("not supported file:{%s}", filePath)
	}
}
