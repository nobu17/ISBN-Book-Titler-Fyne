package zbar

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type ZBarReader struct {
	toolPath string
}

func NewZBarReader(toolPath string) *ZBarReader {
	return &ZBarReader{toolPath}
}

var isbnHeader = "ISBN-13:"

func (z *ZBarReader) GetIsbnFromImage(imagePath string) (string, error) {
	// fmt.Println("command:::", z.toolPath + " -Sisbn13.enable " + imagePath)	
	out, err := exec.Command(z.toolPath, "-Sisbn13.enable", imagePath).CombinedOutput()
	if err != nil {
		// fmt.Println(err)
		return "", err
	}
	isbn := ""
	for i, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(out), -1) {
		fmt.Println(i+1, ":", v)
		if strings.HasPrefix(v, isbnHeader) {
			isbn = strings.Replace(v, isbnHeader, "", -1)
			break
		}
	}
	if isbn == "" {
		return isbn, errors.New("scanning ISBN code is failed")
	}
	return isbn, nil
}
