package settings

import (
	"errors"
	"isbnbook/app/log"
	"strconv"
	"strings"
)

const AppSettingPath = "./app_setting.json"

type AppSettings struct {
	GSPath           string `json:"gsPath"`
	ZBarPath         string `json:"zbarPath"`
	ExtractPages     string `json:"extractPages"`
	RenameOption     string `json:"renameOption"`
	BookReader       string `json:"bookReader"`
	RakutenApiKey    string `json:"rakutenApiKey"`
	AmazonPASettings AmazonPASettings
	json             FileStore
	logger           log.AppLogger
}

type AmazonPASettings struct {
	AssociateId string `json:"associateId"`
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
}

func NewAppSetings() *AppSettings {
	return &AppSettings{
		GSPath:           "gs",
		ZBarPath:         "zbarimg",
		ExtractPages:     "5",
		RenameOption:     Copy.String(),
		BookReader:       OpenBD.String(),
		RakutenApiKey:    "",
		AmazonPASettings: AmazonPASettings{"", "", ""},
		json:             NewJsonSettings(AppSettingPath),
		logger:           log.GetLogger(),
	}
}

func NewAppSetingsWithParam(fs FileStore, logger log.AppLogger) *AppSettings {
	return &AppSettings{
		GSPath:           "gs",
		ZBarPath:         "zbarimg",
		ExtractPages:     "5",
		RenameOption:     Copy.String(),
		BookReader:       OpenBD.String(),
		RakutenApiKey:    "",
		AmazonPASettings: AmazonPASettings{"", "", ""},
		json:             fs,
		logger:           logger,
	}
}

type BookInfoReaderType int

const (
	OpenBD BookInfoReaderType = iota
	NationalLib
	RakutenBook
	AmazonPA
)

func (b BookInfoReaderType) String() string {
	switch b {
	case OpenBD:
		return "OpenBD"
	case NationalLib:
		return "国会図書館"
	case RakutenBook:
		return "楽天ブックAPI"
	case AmazonPA:
		return "AmazonPA API"
	default:
		return "Unknown"
	}
}

type RenameOption int

const (
	Copy RenameOption = iota
	Rename
)

func (b RenameOption) String() string {
	switch b {
	case Copy:
		return "コピーしてリネーム"
	case Rename:
		return "既存ファイルをリネーム"
	default:
		return "Unknown"
	}
}

func (a *AppSettings) Init() {
	loaded := NewAppSetings()
	if err := a.json.Load(loaded); err != nil {
		a.logger.Error("load json is failed.", err)
		// works as default value
		return
	}
	a.GSPath = loaded.GSPath
	a.ZBarPath = loaded.ZBarPath
	a.ExtractPages = loaded.ExtractPages
	a.RenameOption = loaded.RenameOption
	a.BookReader = loaded.BookReader
	a.RakutenApiKey = loaded.RakutenApiKey
	a.AmazonPASettings = loaded.AmazonPASettings
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

	if a.BookReader == RakutenBook.String() {
		if strings.TrimSpace(a.RakutenApiKey) == "" {
			msg = "rakuten API is needed apikey"
			return errors.New(msg)
		}
	}

	if a.BookReader == AmazonPA.String() {
		if strings.TrimSpace(a.AmazonPASettings.AssociateId) == "" {
			msg = "AmazonPA API is needed AssociateId"
			return errors.New(msg)
		}
		if strings.TrimSpace(a.AmazonPASettings.AccessKey) == "" {
			msg = "AmazonPA API is needed AccessKey"
			return errors.New(msg)
		}
		if strings.TrimSpace(a.AmazonPASettings.SecretKey) == "" {
			msg = "AmazonPA API is needed SecretKey"
			return errors.New(msg)
		}
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

func (a *AppSettings) GetSelectableRenames() []string {
	return []string{Copy.String(), Rename.String()}
}

func (a *AppSettings) GetSelectableReader() []string {
	return []string{OpenBD.String(), NationalLib.String(), RakutenBook.String(), AmazonPA.String()}
}
