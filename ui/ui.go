package ui

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/bpicode/twatch/test"
	ui "github.com/gizak/termui"
)

type TwatchUi struct {
	sync.Mutex
	PackageList *PackageList
	TestOutput  *TestOutput
}

var logChan = make(chan ui.Event, 10)

func (twatchUi *TwatchUi) Update(p test.Package) {
	twatchUi.Lock()
	defer twatchUi.Unlock()
	twatchUi.PackageList.Update(p)
	ui.Render(twatchUi.PackageList)
	for _, r := range p.Fails {
		twatchUi.TestOutput.Update(r)
		ui.Render(twatchUi.TestOutput)
		break
	}
}

func InitUi() *TwatchUi {
	err := initTermUi()
	if err != nil {
		log.Fatal(err)
	}

	ui.Merge("twatchLog", logChan)

	ll := NewLogWatcher(logChan)
	pl := NewPackageList()
	to := NewTestOutput()
	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(12, 0, ll)),
		ui.NewRow(ui.NewCol(6, 0, pl), ui.NewCol(6, 0, to)),
	)
	ui.Body.Align()

	twatchUi := &TwatchUi{
		PackageList: pl,
		TestOutput:  to,
	}
	twatchUi.handleClear()

	pressQToQuit()
	handleResize()

	return twatchUi
}

func initTermUi() error {
	return ui.Init()
}

func (twatchUi *TwatchUi) handleClear() {
	ui.Handle("/sys/kbd/c", func(ui.Event) {
		twatchUi.Clear()
	})
}

func (twatchUi *TwatchUi) Clear() {
	twatchUi.Lock()
	defer twatchUi.Unlock()
	ui.Clear()
	twatchUi.PackageList.Clear()
	ui.Render(twatchUi.PackageList)
	twatchUi.TestOutput.Clear()
	ui.Render(twatchUi.TestOutput)
}

func handleResize() {
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Align()
		ui.Render(ui.Body)
	})
}

func pressQToQuit() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})
}

func Close() {
	ui.Close()
}

func Loop() {
	ui.Loop()
}
