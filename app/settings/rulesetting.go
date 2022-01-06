package settings

import (
	"fmt"
	"strings"

	"isbnbook/app/log"
)

const RuleSettingPath = "./rule_setting.json"

type RuleSettings struct {
	RenameRule string
	json       FileStore
	logger     log.AppLogger
}

func NewRuleSettings() *RuleSettings {
	return &RuleSettings{
		RenameRule: "(@[g])[@[a]]@[t]",
		json: NewJsonSettings(RuleSettingPath),
		logger: log.GetLogger(),
	}
}

func NewRuleSettingsWithParam(fs FileStore, logger log.AppLogger) *RuleSettings {
	return &RuleSettings{
		RenameRule: "(@[g])[@[a]]@[t]",
		json: fs,
		logger: logger,
	}
}

func (r *RuleSettings) Init() {
	loaded := NewRuleSettings()
	if err := r.json.Load(loaded); err != nil {
		r.logger.Error("load json is failed.", err)
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
