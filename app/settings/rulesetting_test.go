package settings

import (
	"fmt"
	"testing"
)

func prepareRuleSettings(renameRule string, err error) *RuleSettings {
	jsonStore := RuleSettings{}
	jsonStore.RenameRule = renameRule
	mockStore := newMockFileStore(&jsonStore, err)

	return NewRuleSettingsWithParam(mockStore, mockLogger)
}

func TestRuleSetting_InitSuccess(t *testing.T) {
	const reanameRule = "(@[g])"
	rule := prepareRuleSettings(reanameRule, nil)
	rule.Init()

	if rule.RenameRule != reanameRule {
		t.Errorf("RenameRule should be:{%s}, actual:{%s}", reanameRule, rule.RenameRule)
	}
}

func TestRuleSetting_InitFailed(t *testing.T) {
	const defaultValue = "(@[g])[@[a]]@[t]"
	rule := prepareRuleSettings("", fmt.Errorf(""))
	rule.Init()

	if rule.RenameRule != defaultValue {
		t.Errorf("If Init failed, RenameRule should be:{%s}, actual:{%s}", defaultValue, rule.RenameRule)
	}
}

func TestRuleSetting_ValidateSuccess(t *testing.T) {
	const reanameRule = "(@[g])"
	rule := prepareRuleSettings(reanameRule, nil)
	rule.Init()
	rule.RenameRule = "(@[g])[test]"

	err := rule.Validate()

	if err != nil {
		t.Errorf("Validate result should be nil, actual:{%s}", err)
	}
}

func TestRuleSetting_ValidateFailed(t *testing.T) {
	const reanameRule = "(@[g])"
	rule := prepareRuleSettings(reanameRule, nil)
	rule.Init()
	rule.RenameRule = ""

	err := rule.Validate()

	if err == nil {
		t.Errorf("Validate result should be return error")
	}
}
