package isbn

import (
	"fmt"
	
	"isbnbook/app/command/zbar"
	"isbnbook/app/utils"
)

type barcodeReader struct {
	*zbar.ZBarReader
}

func NewBarcodeReader(toolPath string) *barcodeReader {
	return &barcodeReader{
		ZBarReader: zbar.NewZBarReader(toolPath),
	}
}

func (b *barcodeReader) GetIsbn(dir string) (string, error) {
	// get files from dir
	files, err := utils.GetFilesFromDir(dir, ".jpg")
	if err != nil {
		return "", err
	}
	isbn := ""
	for _, file := range files {
		isbn, err = b.GetIsbnFromImage(file)
		if err != nil {
			fmt.Println("zbar can not get isbn, try next image.", err)
			continue
		}
		if isbn != "" {
			break
		}
	}
	return isbn, nil
}
