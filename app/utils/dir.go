package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"isbnbook/app/log"

	"github.com/google/uuid"
)

var logger = log.GetLogger()

func GetFilesFromDir(dir string, extension string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if !file.IsDir() && path.Ext(file.Name()) == extension {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}
	return paths, nil
}

func MkUniqueDir(baseDir string) (string, error) {
	err := MkDirIfNotExists(baseDir)
	if err != nil {
		return "", err
	}

	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uu := u.String()
	dir := filepath.Join(baseDir, uu)
	err = MkDirIfNotExists(dir)

	return dir, err
}

func MkDirIfNotExists(path string) error {
	logger.Info(fmt.Sprintf("make dir:%s", path))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("try to make")
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			logger.Error("failed to create", err)
			return err
		}
		err = os.Chmod(path, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
