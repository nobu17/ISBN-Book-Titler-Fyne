package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"isbnbook/app/log"
)

var allowExtensions = []string{".pdf", ".zip"}

func GetFiles() ([]string, []string, error) {
	logger := log.GetLogger()

	files := []string{}
	failed := []string{}
	// get file or dir from drag and drop
	if len(os.Args) < 2 {
		return files, failed, nil
	}
	for _, item := range os.Args[1:] {
		file, err := os.Stat(item)
		if os.IsNotExist(err) {
			logger.Warn(fmt.Sprintf("not exists as file:%s", item))
			failed = append(failed, item)
			continue
		}
		if file.IsDir() {
			// only get 1 layer child items
			childs, err := getDirItems(item)
			if err != nil {
				logger.Error("failed to read dir:", err)
				failed = append(failed, item)
				continue
			}
			files = append(files, childs...)
		} else {
			files = append(files, item)
		}
	}
	// filter by extension
	passed, filtered := getFilteredItems(files)
	return passed, append(failed, filtered...), nil
}

func getDirItems(dir string) ([]string, error) {
	files := []string{}
	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}
	for _, file := range paths {
		if !file.IsDir() {
			files = append(files, filepath.Join(dir, file.Name()))
		}
	}
	return files, nil
}

func getFilteredItems(files []string) ([]string, []string) {
	filtered := []string{}
	failed := []string{}
	for _, file := range files {
		if containsSuffix(strings.ToLower(file), allowExtensions) {
			filtered = append(filtered, file)
		} else {
			failed = append(failed, file)
		}
	}
	return filtered, failed
}

func containsSuffix(target string, lists []string) bool {
	for _, item := range lists {
		if strings.HasSuffix(target, item) {
			return true
		}
	}
	return false
}
