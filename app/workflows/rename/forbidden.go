package rename

import (
	"runtime"
	"strings"
)

func replaceForbiddenName(replaceTarge string) string {
	forbiddens := []string{}
	if runtime.GOOS == "windows" {
		forbiddens = []string{"\"", "<", ">", "|", "¥", "/", ":", "?"}
	} else if runtime.GOOS == "darwin" {
		forbiddens = []string{"/", ":", "."}
	} else if runtime.GOOS == "linux" {
		forbiddens = []string{"/"}
	}
	for _, word := range forbiddens {
		replaceTarge = strings.Replace(replaceTarge, word, "_", -1)
	}
	return replaceTarge
}
