package rename

import (
	"errors"
	"testing"

	"isbnbook/app/settings"
	"isbnbook/app/workflows/book"
)

var createRule = func(renamRule string) *settings.RuleSettings {
	r := settings.NewRuleSettingsWithParam(nil, nil)
	r.RenameRule = renamRule
	return r
}

var createAppSetting = func() *settings.AppSettings {
	a := settings.NewAppSetingsWithParam(nil, nil)
	return a
}

func TestGetReplaceName_success(t *testing.T) {
	type input struct {
		rule     *settings.RuleSettings
		book     *book.BookInfo
		expected string
	}
	inputs := []input{
		{createRule("(@[t])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(タイトル)"},
		{createRule("(@[a])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(著者1,著者2)"},
		{createRule("(@[a0])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(著者1)"},
		{createRule("(@[a1])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(著者2)"},
		{createRule("(@[p])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(出版社A)"},
		{createRule("(@[d])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(2021-12)"},
		{createRule("(@[k])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(単行本)"},
		{createRule("(@[g])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(技術書)"},
		{createRule("(@[t])(@[a])(@[a0])(@[a1])(@[p])(@[d])(@[k])(@[g])"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(タイトル)(著者1,著者2)(著者1)(著者2)(出版社A)(2021-12)(単行本)(技術書)"},
		{createRule("(そのまま)"), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(そのまま)"},
		{createRule("(@[t])"), book.NewBookInfo("タ:イトル/", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書"), "(タ_イトル_)"},
	}

	man := newRenameManager()

	for _, d := range inputs {
		actual, err := man.GetReplaceName(d.rule, d.book)
		if err != nil {
			t.Errorf("GetReplaceName result should not have error:%s", err)
		}
		if d.expected != actual {
			t.Errorf("GetReplaceName result should be:%s, actual:%s", d.expected, actual)
		}
	}
}

func TestGetReplaceName_failed(t *testing.T) {
	type input struct {
		rule *settings.RuleSettings
		book *book.BookInfo
	}

	inputs := []input{
		{createRule(""), book.NewBookInfo("タイトル", []string{"著者1", "著者2"}, "出版社A", "2021-12", "単行本", "技術書")},
	}

	man := newRenameManager()

	for _, d := range inputs {
		_, err := man.GetReplaceName(d.rule, d.book)
		if err == nil {
			t.Errorf("GetReplaceName result should have error:%s", d.rule.RenameRule)
		}
	}
}

type MockRenamer struct {
	err  error
	src  string
	dist string
}

func (m *MockRenamer) Rename(src, dist string) error {
	m.src = src
	m.dist = dist
	if m.err != nil {
		return m.err
	}
	return nil
}

func NewMockRenamer(err error) *MockRenamer {
	return &MockRenamer{err, "", ""}
}

func TestRename_success(t *testing.T) {
	// mock file system IF
	mockRenamer := NewMockRenamer(nil)
	savedGetRenamer := getRenamer
	defer func() { getRenamer = savedGetRenamer }()
	getRenamer = func(appSetting *settings.AppSettings) (fileRenamer, error) {
		return mockRenamer, nil
	}
	savedFileExists := fileExists
	defer func() { fileExists = savedFileExists }()
	fileExists = func(filename string) bool {
		return false
	}

	type input struct {
		oldP         string
		newP         string
		expectedFile string
		renamePath   string
	}

	inputs := []input{
		{"old.pdf", "new", "new.pdf", "new.pdf"},
		{"old.zip", "new", "new.zip", "new.zip"},
		{"old.hoge.zip", "new", "new.zip", "new.zip"},
		{"old", "new", "new", "new"},
		{"dir/old.pdf", "new", "new.pdf", "dir/new.pdf"},
	}

	man := newRenameManager()

	for _, p := range inputs {
		actual, err := man.Rename(p.oldP, p.newP, createAppSetting())
		if err != nil {
			t.Errorf("Rename result should not have error:%s", err)
		}
		if actual != p.expectedFile {
			t.Errorf("Rename result should be:%s, actual:%s", p.expectedFile, actual)
		}
		if mockRenamer.dist != p.renamePath {
			t.Errorf("Rename path should be:%s, actual:%s", p.renamePath, mockRenamer.dist)
		}
	}
}

func TestRename_filexists_failed(t *testing.T) {
	// mock file system IF
	mockRenamer := NewMockRenamer(nil)
	savedGetRenamer := getRenamer
	defer func() { getRenamer = savedGetRenamer }()
	getRenamer = func(appSetting *settings.AppSettings) (fileRenamer, error) {
		return mockRenamer, nil
	}
	savedFileExists := fileExists
	defer func() { fileExists = savedFileExists }()
	fileExists = func(filename string) bool {
		return true
	}

	type input struct {
		oldP string
		newP string
	}

	inputs := []input{
		{"old.pdf", "new"},
	}

	man := newRenameManager()

	for _, p := range inputs {
		_, err := man.Rename(p.oldP, p.newP, createAppSetting())
		if err == nil {
			t.Errorf("Rename result should not have error")
		}
	}
}

func TestRename_renamefile_failed(t *testing.T) {
	// mock file system IF
	mockRenamer := NewMockRenamer(errors.New(""))
	savedGetRenamer := getRenamer
	defer func() { getRenamer = savedGetRenamer }()
	getRenamer = func(appSetting *settings.AppSettings) (fileRenamer, error) {
		return mockRenamer, nil
	}
	savedFileExists := fileExists
	defer func() { fileExists = savedFileExists }()
	fileExists = func(filename string) bool {
		return false
	}

	type input struct {
		oldP string
		newP string
	}

	inputs := []input{
		{"old.pdf", "new"},
	}

	man := newRenameManager()

	for _, p := range inputs {
		_, err := man.Rename(p.oldP, p.newP, createAppSetting())
		if err == nil {
			t.Errorf("Rename result should not have error")
		}
	}
}

func TestGetExplaination(t *testing.T) {
	expected := "@[t]:タイトル\n@[a]:著者(複数時はカンマ区切り)\n@[a(num)]:個別著者 (a0:1人目, a1:2人目...)\n@[d]:出版年月\n@[p]:出版社\n@[k]:出版種類(単行本など)\n@[g]:ジャンル"
	man := newRenameManager()

	actual := man.GetExplaination()
	if actual != expected {
		t.Errorf("GetExplaination result should be:%s, actual:%s", expected, actual)
	}
}
