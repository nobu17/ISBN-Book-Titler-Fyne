package isbn

import (
	"strconv"

	"isbnbook/app/command/gs"
)

type pdfExtractor struct {
	filePath   string
	pageOffset int
	*gs.GSPdfReader
}

func newPDFExtractor(toolPath string, filePath string, pageOffset int) *pdfExtractor {
	return &pdfExtractor{
		filePath:    filePath,
		pageOffset:  pageOffset,
		GSPdfReader: gs.NewGSPdfReader(toolPath),
	}
}

func (p pdfExtractor) Extract(workDir string) error {
	page, err := p.GetPageCount(p.filePath)
	if err != nil {
		return err
	}
	// get extract page scopes
	pRanges := GetPageRange(page, p.pageOffset)
	// extract images
	for i, pages := range pRanges {
		err = p.ExtractImageFiles(p.filePath, workDir, pages.start, pages.end, strconv.Itoa(i)+"_")
		if err != nil {
			return err
		}
	}
	return nil
}
