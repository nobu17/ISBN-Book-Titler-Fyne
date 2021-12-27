package isbn

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
)

type zipExtractor struct {
	filePath   string
	pageOffset int
}

func newZipExtractor(filePath string, pageOffset int) *zipExtractor {
	return &zipExtractor{
		filePath, pageOffset,
	}
}

func (z zipExtractor) Extract(workDir string) error {
	r, err := zip.OpenReader(z.filePath)
	if err != nil {
		return err
	}
	defer r.Close()

	files := []zip.File{}
	for _, f := range r.File {
		if !f.FileInfo().IsDir() {
			files = append(files, *f)
		}
	}
	// sort by file names
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
	// get extract page scopes
	pRanges := GetPageRange(len(files), z.pageOffset)
	for _, pRange := range pRanges {
		// extract zip files
		for _, f := range files[pRange.start-1 : pRange.end-1] {
			fmt.Println("file:", f.Name)
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("zip open error:%s", err)
			}
			defer rc.Close()

			buf := make([]byte, f.UncompressedSize)
			_, err = io.ReadFull(rc, buf)
			if err != nil {
				return fmt.Errorf("zip read error:%s", err)
			}
			extractPath := filepath.Join(workDir, f.FileInfo().Name())
			if err = ioutil.WriteFile(extractPath, buf, f.Mode()); err != nil {
				fmt.Println("write error", err)
				return fmt.Errorf("write error:%s", err)
			}
		}
	}

	return nil
}
