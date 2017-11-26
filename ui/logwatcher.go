package ui

import (
	"time"

	log "github.com/Sirupsen/logrus"
	ui "github.com/gizak/termui"
)

func NewLogWatcher(ec chan ui.Event) *ui.Par {
	logLine := ui.NewPar("")
	logLine.Border = false
	logLine.Height = 1
	logLine.TextFgColor = ui.ColorCyan
	log.SetOutput(&writeEventChan{ec: ec})

	logLine.Handle("/twatch/log", func(e ui.Event) {
		if l, ok := e.Data.(string); ok {
			logLine.Text = l
			ui.Render(logLine)
		}
	})
	return logLine
}

type writeEventChan struct {
	ec chan ui.Event
}

func (w *writeEventChan) Write(p []byte) (n int, err error) {
	s := string(p)
	w.ec <- ui.Event{Path: "/twatch/log", Type: "logStatement", Data: s, Time: time.Now().Unix()}
	return len(p), nil
}
