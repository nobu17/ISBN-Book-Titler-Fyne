package utils

import (
	"fmt"
	"strconv"
)

func ConvertToISBN10(isbn13 string) (string, error) {
	base := isbn13[3 : len(isbn13)-1]
	if len(base) != 9 {
		fmt.Println("base:", base)
		return "", fmt.Errorf("incorrect isbn13:%s", isbn13)
	}
	// cal check digit
	total := 0
	for i, data := range base {
		num, err := strconv.Atoi(string(data))
		if err != nil {
			return "", err
		}
		total += (num * (10 - i))
	}
	checkSum := 11 - (total % 11)
	return base + strconv.Itoa(checkSum), nil
}
