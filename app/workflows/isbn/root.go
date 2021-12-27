package isbn

import (
	"fmt"
	"os"

	"isbnbook/app/settings"
	"isbnbook/app/utils"
)

const workDir = "./temp_work"

type IsbnGetWorkFlow struct {
	setting *settings.AppSettings
}

func NewIsbnGetWorkFlow(setting *settings.AppSettings) *IsbnGetWorkFlow {
	return &IsbnGetWorkFlow{
		setting,
	}
}

func (i *IsbnGetWorkFlow) GetIsbn(path string, pageOffset int) (string, error) {
	// prepare work dir
	tempDir, err := utils.MkUniqueDir(workDir)
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir:%s", err)
	}
	defer os.RemoveAll(tempDir)

	extractor, err := NewFileExtractor(i.setting, path)
	if err != nil {
		return "", fmt.Errorf("failed to create extractor:%s", err)
	}
	err = extractor.Extract(tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to extract files:%s", err)
	}
	// get files from dir
	barcode := NewBarcodeReader(i.setting.ZBarPath)
	isbn, err := barcode.GetIsbn(tempDir)
	if err != nil {
		return "", err
	}
	if isbn == "" {
		return isbn, fmt.Errorf("failed to get barcode")
	}
	return isbn, err
}
