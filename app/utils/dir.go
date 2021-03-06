package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

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
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("try to make")
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create. err:%s", err)
		}
		err = os.Chmod(path, 0777)
		if err != nil {
			return fmt.Errorf("failed to change permission. err:%s", err)
		}
	}
	return nil
}

func ChangeWorkDir() error {
	// if on go run command, not change dir
	if RunningThroughGoRun() {
		return nil
	}
	exPath, err := GetAppDir()
	if err != nil {
		return err
	}

	err = os.Chdir(exPath)
	if err != nil {
		return fmt.Errorf("failed to change work dir err:%s", err)
	}
	return nil
}

func GetAppDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get exeutable's dir err:%s", err)
	}
	return filepath.Dir(ex), nil
}

func RunningThroughGoRun() bool {
	executable, err := os.Executable()
	if err != nil {
		return false
	}

	goTmpDir := os.Getenv("GOTMPDIR")
	if goTmpDir != "" {
		return strings.HasPrefix(executable, goTmpDir)
	}

	return strings.HasPrefix(executable, os.TempDir())
}
