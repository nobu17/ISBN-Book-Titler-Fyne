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

const isbnHeader = "ISBN-13:"
const eanCheckHeader = "EAN-13:491" // magazine code
const eanHeader = "EAN-13:"

func (z *ZBarReader) GetIsbnFromImage(imagePath string) (string, error) {
	// fmt.Println("command:::", z.toolPath + " -Sisbn13.enable " + imagePath)	
	out, err := exec.Command(z.toolPath, "-Sisbn13.enable", imagePath).CombinedOutput()
	if err != nil {
		// fmt.Println(err)
		return "", err
	}
	code := ""
	for i, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(string(out), -1) {
		fmt.Println(i+1, ":", v)
		if strings.HasPrefix(v, isbnHeader) {
			code = strings.Replace(v, isbnHeader, "", -1)
			break
		}
		if strings.HasPrefix(v, eanCheckHeader) {
			code = strings.Replace(v, eanHeader, "", -1)
			// try check to isbn is exists
		}
	}
	if code == "" {
		return code, errors.New("scanning ISBN(or EAN 491) code is failed")
	}
	return code, nil
}
