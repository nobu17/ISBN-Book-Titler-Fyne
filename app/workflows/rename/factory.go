package rename

import (
	"fmt"
	"io"
	"os"

	"isbnbook/app/settings"
)

type fileRenamer interface {
	Rename(src, dist string) error
}

var getRenamer = func(appSetting *settings.AppSettings) (fileRenamer, error) {
	if appSetting.RenameOption == settings.Copy.String() {
		return &copyRenmaer{}, nil
	} else if appSetting.RenameOption == settings.Rename.String() {
		return &renmaer{}, nil
	} else {
		return nil, fmt.Errorf("cant not get fileRenamer")
	}
}

type renmaer struct {
}

func (r *renmaer) Rename(src, dist string) error {
	return os.Rename(src, dist)
}

type copyRenmaer struct {
}

func (c *copyRenmaer) Rename(src, dist string) error {
	srcfile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s, error:%s", src, err.Error())
	}
	defer srcfile.Close()

	dst, err := os.Create(dist)
	if err != nil {
		return fmt.Errorf("failed to crate dist file %s, error:%s", dist, err.Error())
	}
	defer dst.Close()

	_, err = io.Copy(dst, srcfile)
	if err != nil {
		return fmt.Errorf("failed to copy:%s, dist:%s, error:%s", src, dist, err.Error())
	}
	return nil
}
