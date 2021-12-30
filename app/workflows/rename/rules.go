package rename

import (
	"fmt"
	"strings"

	"isbnbook/app/workflows/book"
)

type RenameRule interface {
	GetReplacedName(bookInfo *book.BookInfo, baseStr string) string
	GetExplaination() string
}

type replaceLogic struct {
	word         string
	replaceParam string
	explanation  string
}

func newReplaceLogic(word string, explanation string) *replaceLogic {
	replaceWord := "@[" + word + "]"
	explanationMessage := replaceWord + ":" + explanation
	return &replaceLogic{
		word,
		replaceWord,
		explanationMessage,
	}
}

func (r *replaceLogic) Replace(replaceTarge string, newStr string) string {
	return strings.Replace(replaceTarge, r.replaceParam, newStr, -1)
}

func (r *replaceLogic) GetExplaination() string {
	return r.explanation
}

type titleRule struct {
	*replaceLogic
}

func newTitleRule() *titleRule {
	return &titleRule{
		replaceLogic: newReplaceLogic("t", "タイトル"),
	}
}

func (t *titleRule) GetReplacedName(bookInfo *book.BookInfo, baseStr string) string {
	return t.Replace(baseStr, bookInfo.Title)
}

type publisherRule struct {
	*replaceLogic
}

func newPublisherRule() *publisherRule {
	return &publisherRule{
		replaceLogic: newReplaceLogic("p", "出版社"),
	}
}

func (t *publisherRule) GetReplacedName(bookInfo *book.BookInfo, baseStr string) string {
	return t.Replace(baseStr, bookInfo.Publisher)
}

type authorsRule struct {
	*replaceLogic
}

func newAuthorRule() *authorsRule {
	return &authorsRule{
		replaceLogic: newReplaceLogic("a", "著者(複数時はカンマ区切り)"),
	}
}

func (t *authorsRule) GetReplacedName(bookInfo *book.BookInfo, baseStr string) string {
	return t.Replace(baseStr, strings.Join(bookInfo.Authors, ","))
}

type separateAuthorRule struct {
}

func newSeparateAuthorRule() *separateAuthorRule {
	return &separateAuthorRule{}
}

func (t *separateAuthorRule) GetReplacedName(bookInfo *book.BookInfo, baseStr string) string {
	replaced := baseStr
	for i, author := range bookInfo.Authors {
		replace := newReplaceLogic(fmt.Sprintf("a%d", i), "")
		replaced = replace.Replace(replaced, author)
	}
	return replaced
}

func (r *separateAuthorRule) GetExplaination() string {
	return "@[a(num)]:個別著者 (a0:1人目, a1:2人目...)"
}

type dateRule struct {
	*replaceLogic
}

func newDateRule() *dateRule {
	return &dateRule{
		replaceLogic: newReplaceLogic("d", "出版年月"),
	}
}

func (d *dateRule) GetReplacedName(bookInfo *book.BookInfo, baseStr string) string {
	return d.Replace(baseStr, bookInfo.Date)
}

type genreRule struct {
	*replaceLogic
}

func newGenreRule() *genreRule {
	return &genreRule{
		replaceLogic: newReplaceLogic("g", "ジャンル"),
	}
}

func (t *genreRule) GetReplacedName(bookInfo *book.BookInfo, baseStr string) string {
	return t.Replace(baseStr, bookInfo.Genre)
}

type kindRule struct {
	*replaceLogic
}

func newKindRule() *kindRule {
	return &kindRule{
		replaceLogic: newReplaceLogic("k", "出版種類(単行本など)"),
	}
}

func (t *kindRule) GetReplacedName(bookInfo *book.BookInfo, baseStr string) string {
	return t.Replace(baseStr, bookInfo.Kind)
}
