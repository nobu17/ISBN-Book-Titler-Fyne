package isbn

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/gen2brain/go-unarr"
)

type rarExtractor struct {
	filePath   string
	pageOffset int
}

func newRarExtractor(filePath string, pageOffset int) *rarExtractor {
	return &rarExtractor{
		filePath, pageOffset,
	}
}

func (z rarExtractor) Extract(workDir string) error {
	a, err := unarr.NewArchive(z.filePath)
	if err != nil {
		return fmt.Errorf("failed to extract rar:", err)
	}
	defer a.Close()

	list, err := a.List()
	if err != nil {
		return fmt.Errorf("failed to read file list from rar:", err)
	}
	// sort by file names
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	// get extract page scopes
	pRanges := GetPageRange(len(list), z.pageOffset)
	for _, pRange := range pRanges {
		// extract zip files
		for _, f := range list[pRange.start-1 : pRange.end-1] {
			extractPath := filepath.Join(workDir, filepath.Base(f))
			err := a.EntryFor(f);
			if err != nil {
				return fmt.Errorf("failed to read entry:%s", err)		
			}
			data, err := a.ReadAll()
			if err != nil {
				return fmt.Errorf("failed to read file:%s", err)		
			}
			if err = ioutil.WriteFile(extractPath, data, os.ModePerm); err != nil {
				fmt.Println("write error", err)
				return fmt.Errorf("write error:%s", err)
			}
		}
	}
	return nil
}
