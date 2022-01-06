package rename

import (
	"testing"

	"isbnbook/app/settings"
	"isbnbook/app/workflows/book"
)

var createRule = func(renamRule string) *settings.RuleSettings {
	r := settings.NewRuleSettingsWithParam(nil, nil)
	r.RenameRule = renamRule
	return r
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
		rule     *settings.RuleSettings
		book     *book.BookInfo
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

func TestGetExplaination(t *testing.T) {
	expected := "@[t]:タイトル\n@[a]:著者(複数時はカンマ区切り)\n@[a(num)]:個別著者 (a0:1人目, a1:2人目...)\n@[d]:出版年月\n@[p]:出版社\n@[k]:出版種類(単行本など)\n@[g]:ジャンル" 
	man := newRenameManager()

	actual := man.GetExplaination()
	if actual != expected {
		t.Errorf("GetExplaination result should be:%s, actual:%s", expected, actual)
	}
}
