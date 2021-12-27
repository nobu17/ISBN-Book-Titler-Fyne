package settings

import (
	"fmt"
	"strings"

	"isbnbook/app/log"
)

var logger = log.GetLogger()

const RuleSettingPath = "./rule_setting.json"

type RuleSettings struct {
	RenameRule string
	json       *JsonSettings
}

func NewRuleSettings() *RuleSettings {
	return &RuleSettings{
		RenameRule: "(@[g])[@[a]]@[t]",
		json: NewJsonSettings(RuleSettingPath),
	}
}

func (r *RuleSettings) Init() {
	loaded := RuleSettings{}
	if err := r.json.Load(&loaded); err != nil {
		logger.Error("load json is failed.", err)
		return
	}
	r.RenameRule = loaded.RenameRule
}

func (r *RuleSettings) SaveSetting() error {
	return r.json.Save(&r)
}

func (r *RuleSettings) Validate() error {
	if strings.Trim(r.RenameRule, "") == "" {
		return fmt.Errorf("rename rule should not be empty")
	}
	return nil
}
