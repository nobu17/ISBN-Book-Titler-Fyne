package settings

import (
	"errors"
	"strconv"
	"strings"
)

const AppSettingPath = "./app_setting.json"

type AppSettings struct {
	GSPath       string `json:"gsPath"`
	ZBarPath     string `json:"zbarPath"`
	ExtractPages string `json:"extractPages"`
	BookReader   string `json:"bookReader"`
	json         *JsonSettings
}

func NewAppSetings() *AppSettings {
	return &AppSettings{
		GSPath:       "gs",
		ZBarPath:     "zbarimg",
		ExtractPages: "5",
		BookReader:   OpenBD.String(),
		json:         NewJsonSettings(AppSettingPath),
	}
}

type BookInfoReaderType int

const (
	OpenBD BookInfoReaderType = iota
	NationalLib
)

func (b BookInfoReaderType) String() string {
	switch b {
	case OpenBD:
		return "OpenBD"
	case NationalLib:
		return "国会図書館"
	default:
		return "Unknown"
	}
}

func (a *AppSettings) Init() {
	loaded := AppSettings{}
	if err := a.json.Load(&loaded); err != nil {
		logger.Error("load json is failed.", err)
		// works as default value
		return
	}
	a.GSPath = loaded.GSPath
	a.ZBarPath = loaded.ZBarPath
	a.ExtractPages = loaded.ExtractPages
	a.BookReader = loaded.BookReader
}

func (a *AppSettings) SaveSetting() error {
	return a.json.Save(&a)
}

func (a *AppSettings) Validate() error {
	isValid := false
	msg := ""
	pages := a.GetSelectablePages()
	for _, item := range pages {
		if strings.TrimSpace(item) == strings.TrimSpace(a.ExtractPages) {
			isValid = true
			break
		}
	}
	if !isValid {
		msg = "selectable page is not in range"
		return errors.New(msg)
	}

	isValid = false
	readers := a.GetSelectableReader()
	for _, item := range readers {
		if strings.TrimSpace(item) == strings.TrimSpace(a.BookReader) {
			isValid = true
			break
		}
	}
	if !isValid {
		msg = "selectable reader is not in range"
		return errors.New(msg)
	}

	return nil
}

const DefaultPageInt = 5

func (a *AppSettings) GetPagesInt() int {
	num, err := strconv.Atoi(a.ExtractPages)
	if err != nil {
		return DefaultPageInt
	}
	return num
}

func (a *AppSettings) GetSelectablePages() []string {
	return []string{"1", "2", "3", "4", "5", "6", "7"}
}

func (a *AppSettings) GetSelectableReader() []string {
	return []string{OpenBD.String(), NationalLib.String()}
}
