package settings

import (

	//"reflect"
	"testing"

	"isbnbook/app/log"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	mockLogger = log.NewMockLogger()
}

func teardown() {
}


type mockFileStore struct {
	Target interface{}
	Err    error
}

func newMockFileStore(target interface{}, err error) *mockFileStore {
	return &mockFileStore{Target: target, Err: err}
}

func (m *mockFileStore) Load(loadTarget interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	switch value := loadTarget.(type) {
	case *RuleSettings:
		rule, ok := m.Target.(*RuleSettings)
		if ok {
			value.RenameRule = rule.RenameRule
		} else {
			panic("can not convert *RuleSettings")
		}
	case *AppSettings:
		app, ok := m.Target.(*AppSettings)
		if ok {
			value.GSPath = app.GSPath
			value.ZBarPath = app.ZBarPath
			value.BookReader = app.BookReader
			value.ExtractPages = app.ExtractPages
			value.RenameOption = app.RenameOption
			value.RakutenApiKey = app.RakutenApiKey
			value.AmazonPASettings = app.AmazonPASettings
		} else {
			panic("can not convert *AppSettings")
		}
	case **AppSettings:
		panic("pooooooo")
	default:
		panic("can not convert any types")
	}

	return nil
}
func (m *mockFileStore) Save(saveTarge interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	m.Target = saveTarge
	return nil
}

var mockLogger log.AppLogger
