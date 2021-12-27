package book

import (
	"fmt"
	
	"isbnbook/app/settings"
)

type BookInfo struct {
	Title     string
	Authors   []string
	Publisher string
	Date      string
	Kind      string
	Genre     string
}

type BookApi interface {
	GetBookInfo(isbn13 string) (*BookInfo, error)
}

func NewBookInfo(titile string, authors []string, publisher string, date string, kind string, genre string) *BookInfo {
	return &BookInfo{
		Title:     titile,
		Authors:   authors,
		Publisher: publisher,
		Date:      date,
		Kind:      kind,
		Genre:     genre,
	}
}

func GetBookInfo(isbn13 string, settings *settings.AppSettings) (*BookInfo, error) {
	api, err := getBookApi(settings)
	if err != nil {
		return nil, err
	}
	return api.GetBookInfo(isbn13)
}

func getBookApi(setting *settings.AppSettings) (BookApi, error) {
	var api BookApi
	switch setting.BookReader {
	case settings.OpenBD.String():
		api = NewOpenBDReader()
	case settings.NationalLib.String():
		api = NewNationalLibReader()
	}
	if api != nil {
		return api, nil
	}
	return nil, fmt.Errorf("failed get api")
}
