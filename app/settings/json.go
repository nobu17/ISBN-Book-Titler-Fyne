package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type JsonSettings struct {
	Path string
}

func NewJsonSettings(path string) *JsonSettings {
	return &JsonSettings{
		path,
	}
}

func (j *JsonSettings) Load(loadTarget interface{}) error {
	data, err := ioutil.ReadFile(j.Path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &loadTarget); err != nil {
		return err
	}
	return nil
}

func (j *JsonSettings) Save(saveTarge interface{}) error {
	json, err := json.Marshal(saveTarge)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(j.Path, json, os.ModePerm); err != nil {
		return nil
	}
	return nil
}
