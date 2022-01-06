package settings

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBookInfoReaderType_String(t *testing.T) {
	tables := map[BookInfoReaderType]string{
		OpenBD:      "OpenBD",
		NationalLib: "国会図書館",
		RakutenBook: "楽天ブックAPI",
		AmazonPA:    "AmazonPA API",
		4:           "Unknown",
	}
	for k, v := range tables {
		if k.String() != v {
			t.Errorf("String is expected:%s, actual;%s", v, k.String())
		}
	}
}

func prepareAppSettings(mockData AppSettings, err error) *AppSettings {
	mockStore := newMockFileStore(&mockData, err)

	return NewAppSetingsWithParam(mockStore, mockLogger)
}

func assertAppSettings(expected, actual *AppSettings, t *testing.T) {
	if actual.GSPath != expected.GSPath {
		t.Errorf("GSPath should be:{%s}, actual:{%s}", expected.GSPath, actual.GSPath)
	}
	if actual.ZBarPath != expected.ZBarPath {
		t.Errorf("GSPath should be:{%s}, actual:{%s}", expected.ZBarPath, actual.ZBarPath)
	}
	if actual.BookReader != expected.BookReader {
		t.Errorf("BookReader should be:{%s}, actual:{%s}", expected.BookReader, actual.BookReader)
	}
	if actual.ExtractPages != expected.ExtractPages {
		t.Errorf("ExtractPages should be:{%s}, actual:{%s}", expected.ExtractPages, actual.ExtractPages)
	}
	if actual.RakutenApiKey != expected.RakutenApiKey {
		t.Errorf("RakutenApiKey should be:{%s}, actual:{%s}", expected.RakutenApiKey, actual.RakutenApiKey)
	}
	if !reflect.DeepEqual(actual.AmazonPASettings, expected.AmazonPASettings) {
		t.Errorf("AmazonPASettings should be:{%s}, actual:{%s}", expected.AmazonPASettings, actual.AmazonPASettings)
	}
}

func TestAppSettings_Init_Success(t *testing.T) {
	mockData := AppSettings{}
	mockData.GSPath = "gsp"
	mockData.ZBarPath = "zbbb"
	mockData.BookReader = "reader"
	mockData.ExtractPages = "1"
	mockData.RakutenApiKey = "key"
	mockData.AmazonPASettings = AmazonPASettings{"Id", "Key", "SecKey"}

	app := prepareAppSettings(mockData, nil)
	app.Init()

	assertAppSettings(&mockData, app, t)
}

func TestAppSettings_Init_Failed(t *testing.T) {
	defaultVal := AppSettings{
		GSPath:           "gs",
		ZBarPath:         "zbarimg",
		ExtractPages:     "5",
		BookReader:       OpenBD.String(),
		RakutenApiKey:    "",
		AmazonPASettings: AmazonPASettings{"", "", ""},
		json:             nil,
		logger:           nil,
	}

	app := prepareAppSettings(AppSettings{}, fmt.Errorf(""))
	app.Init()

	assertAppSettings(&defaultVal, app, t)
}

func TestAppSettings_Validate_BookReader_Success(t *testing.T) {
	tables := []BookInfoReaderType{
		OpenBD,
		NationalLib,
		RakutenBook,
		AmazonPA,
	}
	for _, reader := range tables {
		app := AppSettings{
			GSPath:           "gs",
			ZBarPath:         "zbarimg",
			ExtractPages:     "5",
			BookReader:       reader.String(),
			RakutenApiKey:    "key",
			AmazonPASettings: AmazonPASettings{"Id", "Key", "SecKey"},
			json:             nil,
			logger:           nil,
		}

		if err := app.Validate(); err != nil {
			t.Errorf("should be no error at BookReader:{%s}, err:{%s}", app.BookReader, err)
		}
	}
}

func TestAppSettings_Validate_BookReader_Failed(t *testing.T) {
	app := AppSettings{
		GSPath:           "gs",
		ZBarPath:         "zbarimg",
		ExtractPages:     "5",
		BookReader:       "err",
		RakutenApiKey:    "key",
		AmazonPASettings: AmazonPASettings{"Id", "Key", "SecKey"},
		json:             nil,
		logger:           nil,
	}

	if err := app.Validate(); err == nil {
		t.Errorf("should hasve error when incorrect BookReader:%s", app.BookReader)
	}
}

func TestAppSettings_Validate_BookReader_RakutenApi_Failed(t *testing.T) {
	app := AppSettings{
		GSPath:           "gs",
		ZBarPath:         "zbarimg",
		ExtractPages:     "5",
		BookReader:       RakutenBook.String(),
		RakutenApiKey:    "",
		AmazonPASettings: AmazonPASettings{"Id", "Key", "SecKey"},
		json:             nil,
		logger:           nil,
	}

	if err := app.Validate(); err == nil {
		t.Errorf("should hasve error when RakutenApiKey is empty")
	}
}

func TestAppSettings_Validate_BookReader_AmazonPA_Failed(t *testing.T) {
	tables := []AmazonPASettings{
		{"", "Key", "SecKey"},
		{" ", "Key", "SecKey"},
		{"Id", "", "SecKey"},
		{"Id", " ", "SecKey"},
		{"Id", "Key", ""},
		{"Id", "Key", " "},
		{"", "", ""},
		{" ", " ", " "},
	}
	for _, pa := range tables {
		app := AppSettings{
			GSPath:           "gs",
			ZBarPath:         "zbarimg",
			ExtractPages:     "5",
			BookReader:       AmazonPA.String(),
			RakutenApiKey:    "",
			AmazonPASettings: pa,
			json:             nil,
			logger:           nil,
		}
	
		if err := app.Validate(); err == nil {
			t.Errorf("should hasve error when AmazonPASetting some values are empty:{%s}", pa)
		}
	}
}

func TestAppSettings_GetPagesInt(t *testing.T) {
	tables := map[string]int{
		"1" : 1,
		"5" : 5,
		"a" : 5, //default 
	}
	for k, v := range tables {
		app := AppSettings{
			GSPath:           "gs",
			ZBarPath:         "zbarimg",
			ExtractPages:     k,
			BookReader:       AmazonPA.String(),
			RakutenApiKey:    "",
			AmazonPASettings: AmazonPASettings{"Id", "Key", "SecKey"},
			json:             nil,
			logger:           nil,
		}
	
		if v != app.GetPagesInt() {
			t.Errorf("GetPagesInt should be:%d when ExtractPage is:%s", v, k)
		}
	}
}

func TestAppSettings_GetSelectablePages(t *testing.T) {
	expected := []string{"1", "2", "3", "4", "5", "6", "7"}

	app := AppSettings{
		GSPath:           "gs",
		ZBarPath:         "zbarimg",
		ExtractPages:     "5",
		BookReader:       RakutenBook.String(),
		RakutenApiKey:    "",
		AmazonPASettings: AmazonPASettings{"Id", "Key", "SecKey"},
		json:             nil,
		logger:           nil,
	}

	actual := app.GetSelectablePages()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("GetSelectablePages should be:%s when ExtractPage is:%s", expected, actual)
	}
}

func TestAppSettings_GetSelectableReader(t *testing.T) {
	expected := []string{OpenBD.String(), NationalLib.String(), RakutenBook.String(), AmazonPA.String()}

	app := AppSettings{
		GSPath:           "gs",
		ZBarPath:         "zbarimg",
		ExtractPages:     "5",
		BookReader:       RakutenBook.String(),
		RakutenApiKey:    "",
		AmazonPASettings: AmazonPASettings{"Id", "Key", "SecKey"},
		json:             nil,
		logger:           nil,
	}

	actual := app.GetSelectableReader()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("GetSelectableReader should be:%s when ExtractPage is:%s", expected, actual)
	}
}