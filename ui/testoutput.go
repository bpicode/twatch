package ui

import (
	"sync"

	"github.com/bpicode/twatch/test"
	ui "github.com/gizak/termui"
)

type TestOutput struct {
	sync.Mutex
	*ui.List
}

func (t *TestOutput) Update(r *test.Result) {
	t.Lock()
	defer t.Unlock()
	t.BorderLabel = r.Name
	t.Items = r.Out
	t.Height = len(t.Items) + 2
}

func (t *TestOutput) Clear() {
	t.Lock()
	defer t.Unlock()
	t.BorderLabel = "<test result>"
	t.BorderLabelFg = ui.ColorWhite & ui.AttrBold
	t.Height = 2
}

func NewTestOutput() *TestOutput {
	t := &TestOutput{List: ui.NewList()}
	t.BorderLabel = "<test result>"
	t.BorderLabelFg = ui.ColorWhite & ui.AttrBold
	t.Height = 2
	return t
}

