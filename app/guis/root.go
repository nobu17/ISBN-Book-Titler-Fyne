package guis

import "fyne.io/fyne/v2"

type Displayer interface {
	Display(w *fyne.Window, param interface{}) *fyne.Container
}

type cacheDisplay struct {
	current    *fyne.Container
	createFunc func(w *fyne.Window, param interface{}) *fyne.Container
}

func (c *cacheDisplay) Display(w *fyne.Window, param interface{}) *fyne.Container {
	if c.current == nil {
		c.current = c.createFunc(w, param)
	}

	return c.current
}

func NewSettingDisplay() Displayer {
	return &cacheDisplay{
		current: nil,
		createFunc: func(w *fyne.Window, param interface{}) *fyne.Container {
			return getSettingContent()
		},
	}
}

func NewRuleDisplay() Displayer {
	return &cacheDisplay{
		current: nil,
		createFunc: func(w *fyne.Window, param interface{}) *fyne.Container {
			return getRuleContent()
		},
	}
}

func NewTestDisplay() Displayer {
	return &cacheDisplay{
		current: nil,
		createFunc: func(w *fyne.Window, param interface{}) *fyne.Container {
			return getTestContent(w)
		},
	}
}

func NewResultsDisplay() Displayer {
	return &cacheDisplay{
		current: nil,
		createFunc: func(w *fyne.Window, param interface{}) *fyne.Container {
			f, ok := param.([]string)
			if ok {
				return getResultContent(w, f)
			}
			panic("can not display. parameter is incorrect")
		},
	}
}

func NewVersionDisplay() Displayer {
	return &cacheDisplay{
		current: nil,
		createFunc: func(w *fyne.Window, param interface{}) *fyne.Container {
			return getVersionContent(w)
		},
	}
}
