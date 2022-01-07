package versions

import (
	"fmt"
	"strconv"
	"strings"

	"isbnbook/app/log"
	"isbnbook/app/repos"
)

type Version struct {
	Main   int
	Mainor int
	Patch  int
}

type Client interface {
	GetCurrent() *Version
	GetLatest() (*Version, error)
}

type gitHub struct {
	client repos.Client
	logger log.AppLogger
}

func NewClient() (Client, error) {
	cli, err := repos.NewClient("https://raw.githubusercontent.com/nobu17/ISBN-Book-Titler-Fyne")
	if err != nil {
		return nil, err
	}
	log := log.GetLogger()
	return &gitHub{
		client: cli,
		logger: log,
	}, nil
}

func (g *gitHub) GetCurrent() *Version {
	return getCurrent()
}

func (g *gitHub) GetLatest() (*Version, error) {
	bytes, err := g.client.Get("main/app/versions/const.go", map[string]string{})
	if err != nil {
		g.logger.Error("failed to get from github sorce", err)
		return nil, fmt.Errorf("failed to get version", err)
	}
	str := string(bytes)
	return g.pharse(str)
}

func (g *gitHub) pharse(source string) (*Version, error) {
	const mainPrefix = "currentMain  = "
	const mainorPrefix = "currentSub  = "
	const patchPrefix = "currentPatch = "
	mainNo, mainorNo, patchNo := -1, -1, -1

	strs := strings.Split(source, "\n")
	for _, s := range strs {
		g.tryGet(s, mainPrefix, &mainNo)
		g.tryGet(s, mainorPrefix, &mainorNo)
		g.tryGet(s, patchPrefix, &patchNo)
	}

	ver := Version{mainNo, mainorNo, patchNo}
	if err := ver.Validate(); err != nil {
		return nil, fmt.Errorf("can not correct version from source:%s", err)
	}

	return &ver, nil
}

func (g *gitHub) tryGet(source, prefix string, no *int) {
	if strings.Contains(source, prefix) {
		sp := strings.Split(source, prefix)
		if n, err := strconv.Atoi(sp[len(sp)-1]); err == nil {
			*no = n
		}
	}
}

func (v *Version) Validate() error {
	errMsgs := []string{}
	if v.Main < 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("main is incorrect:%d", v.Main))
	}
	if v.Mainor < 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("mainor is incorrect:%d", v.Mainor))
	}
	if v.Patch < 0 {
		errMsgs = append(errMsgs, fmt.Sprintf("patch is incorrect:%d", v.Patch))
	}
	
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}
