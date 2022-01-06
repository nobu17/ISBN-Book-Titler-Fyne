package utils

import "testing"

func TestConvertToISBN10_Success(t *testing.T) {
	tables := map[string]string {
		"9784286131016": "4286131017", // isbn13, expected isbn10
	}
	for isbn13, isbn10 := range tables {
		actual, err := ConvertToISBN10(isbn13)
		if err != nil {
			t.Errorf("should not be error convert correct isbn13:%s, err:%s", isbn13, err)
		}
		if actual != isbn10 {
			t.Errorf("convert result is not expected isbn13:%s, actual:%s, expected:%s", isbn13, actual, isbn10)	
		}
	}
}

func TestConvertToISBN10_Failed(t *testing.T) {
	tables := []string {
		"97842861310161", // incorrect digits (+1)
		"978428613101", // incorrect digits (-1)
		"978b286131016", // not number
	}
	for _, isbn13 := range tables {
		_, err := ConvertToISBN10(isbn13)
		if err == nil {
			t.Errorf("should be error convert correct isbn13:%s", isbn13)
		}
	}
}