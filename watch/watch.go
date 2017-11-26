package watch

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
)

func Pwd(onRelevantChange func()) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("Cannot watch file system: %v", err)
		return
	}
	err = w.Add(".")
	if err != nil {
		log.Errorf("Cannot watch dir '.': %v", err)
	}
	err = filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !f.IsDir() {
			return nil
		}
		if !shouldWatch(path) {
			return nil
		}
		errAdd := w.Add(path)
		if errAdd != nil {
			log.Errorf("Cannot watch dir '%s' %v", path, errAdd)
		}
		return nil
	})

	for {
		select {
		case event := <-w.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}
			if strings.HasSuffix(event.Name, ".go") {
				log.Infof("Source file event for '%s': %s", event.Name, event.Op)
				go onRelevantChange()
			}
		case err := <-w.Errors:
			log.Errorf("error: %v", err)
		}
	}
}

func shouldWatch(path string) bool {
	return !strings.HasPrefix(path, "vendor") &&
		!strings.HasPrefix(path, "build") &&
		!strings.HasPrefix(path, ".")
}
