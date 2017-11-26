package main

import (
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bpicode/twatch/test"
	"github.com/bpicode/twatch/ui"
	"github.com/bpicode/twatch/watch"
	"github.com/gizak/termui"
)

var running = int64(0)

var pChan = make(chan termui.Event, 10)

func main() {
	twatchUi := ui.InitUi()
	defer ui.Close()

	handleTestResultUpdates(twatchUi)
	handleRetest(twatchUi)

	go run(twatchUi)
	go watch.Pwd(func() { run(twatchUi) })

	ui.Loop()
}

func run(twatchUi *ui.TwatchUi) {
	defer atomic.AddInt64(&running, -1)
	if atomic.AddInt64(&running, 1) != 1 {
		log.Warn("Skipped re-test, because tests are running")
		return
	}
	twatchUi.Clear()
	test.RunProject(func(p test.Package) {
		pChan <- termui.Event{Path:"/twatch/package", Type:"packageDone", Data: p, Time: time.Now().Unix()}
	})
}

func handleRetest(twatchUi *ui.TwatchUi) {
	termui.Handle("/sys/kbd/r", func(_ termui.Event) {
		go run(twatchUi)
	})
}

func handleTestResultUpdates(twatchUi *ui.TwatchUi) {
	termui.Merge("twatchPackage", pChan)
	termui.Handle("/twatch/package", func(e termui.Event) {
		if p, ok := e.Data.(test.Package); ok {
			twatchUi.Update(p)
		}
	})
}
