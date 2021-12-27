package rename

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"isbnbook/app/settings"
	"isbnbook/app/workflows/book"
)

type renameManager struct {
	rules *[]RenameRule
}

func newRenameManager() *renameManager {
	rules := []RenameRule{}
	rules = append(rules, newTitleRule())
	rules = append(rules, newAuthorRule())
	rules = append(rules, newDateRule())
	rules = append(rules, newPublisherRule())
	rules = append(rules, newKindRule())
	rules = append(rules, newGenreRule())
	return &renameManager{&rules}
}

func (r *renameManager) Rename(path string, renameSetting *settings.RuleSettings, bookInfo *book.BookInfo) (string, error) {
	extname := filepath.Ext(path)
	dirname, _ := filepath.Split(path)

	replacedname := renameSetting.RenameRule
	for _, rule := range *r.rules {
		replacedname = rule.GetReplacedName(bookInfo, replacedname)
	}
	if len(strings.Trim(replacedname, "")) <= 0 {
		return "", fmt.Errorf("replace name is empty")
	}
	replacedname += extname
	renamePath := filepath.Join(dirname, replacedname)

	if fileExists(renamePath) {
		return "", fmt.Errorf("already file exists:%s", renamePath)
	}

	if err := os.Rename(path, renamePath); err != nil {
		return "", err
	}
	return replacedname, nil
}

func (r *renameManager) GetExplaination() string {
	explanations := []string{}
	for _, rule := range *r.rules {
		explanations = append(explanations, rule.GetExplaination())
	}
	return strings.Join(explanations, "\n")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
