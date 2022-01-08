package versions

import (
	"reflect"
	"testing"

	"isbnbook/app/log"
	"isbnbook/app/repos"
)

func TestGetLatest_integration(t *testing.T) {
	cli, err := NewClientWithLogger(log.NewMockLogger())
	if err != nil {
		t.Errorf("failed to init:%s", err)
	}
	act, err := cli.GetLatest()
	if err != nil {
		t.Errorf("failed to GetLatest:%s", err)
	}

	exp := cli.GetCurrent()

	if !reflect.DeepEqual(exp, act) {
		t.Errorf("Version is mismatch. expecte:%v, actual:%v", exp, act)
	}
}

func TestGetLatest_success(t *testing.T) {

	tables := []string{
		"currentMain  = 0\ncurrentMinor   = 0\ncurrentPatch = 0",
		"currentMain  = 1\ncurrentMinor   = 2\ncurrentPatch = 3",
	}

	for _, data := range tables {
		mockFunc := func() ([]byte, error) {
			return []byte(data), nil
		}
		cli := NewClientWithParam(repos.NewMockClient(mockFunc), log.NewMockLogger())
		_, err := cli.GetLatest()
		if err != nil {
			t.Errorf("should not have error:%s, input:%s", err, data)
		}
	}
}

func TestGetLatest_failed(t *testing.T) {

	tables := []string{
		"",
		// only 1 case
		"currentMain  = 0",
		"currentMinor   = 0",
		// only 2 cases
		"currentMain  = 0\ncurrentMinor   = 0",
		// iregullar param
		"currentMain  = -1\ncurrentMinor   = 0\ncurrentPatch = 0",
		"currentMain  = 0\ncurrentMinor   = -1\ncurrentPatch = 0",
		"currentMain  = 0\ncurrentMinor   = 0\ncurrentPatch = -1",
		"currentMain  = a\ncurrentMinor   = 0\ncurrentPatch = 0",
	}

	for _, data := range tables {
		mockFunc := func() ([]byte, error) {
			return []byte(data), nil
		}
		cli := NewClientWithParam(repos.NewMockClient(mockFunc), log.NewMockLogger())
		_, err := cli.GetLatest()
		if err == nil {
			t.Errorf("should have error:%s", data)
		}
	}
}
