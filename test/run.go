package test

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Result struct {
	Name     string
	Out      []string
	err      error
	duration time.Duration
}

type Package struct {
	Name     string
	Passes   []*Result
	Fails    []*Result
	duration time.Duration
}

type Project struct {
	Packages []*Package
	duration time.Duration
}

func (p *Package) HasErrors() bool {
	return len(p.Fails) > 0
}

func (p *Package) HasTests() bool {
	return p.NumTests() > 0
}

func (p *Package) NumTests() int {
	return p.NumFails() + p.NumPasses()
}

func (p *Package) NumFails() int {
	return len(p.Fails)
}

func (p *Package) NumPasses() int {
	return len(p.Passes)
}

type Handler func(p Package)

func RunProject(h Handler) {
	rc, doneRun := runTests()
	results, doneEmit := emitEvents(rc)
	doneCollect := collect(results, h)

	<-doneRun
	<-doneEmit
	<-doneCollect
}

func collect(results chan Package, h Handler) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for r := range results {
			h(r)
		}
	}()
	return done
}

func emitEvents(rc io.ReadCloser) (chan Package, chan struct{}) {
	done := make(chan struct{})
	results := make(chan Package, 1)

	go func() {
		defer close(done)
		defer close(results)
		defer rc.Close()
		scanner := bufio.NewScanner(rc)
		currentPackage := new(Package)
		lines := make([]string, 0, 0)
		currentResult := new(Result)
		for scanner.Scan() {
			line := scanner.Text()
			if _, ok := testStarted(line); ok {
				lines = make([]string, 0, 0)
			} else if testName, ok := testSuccess(line); ok {
				currentResult = &Result{Name: testName}
				currentPackage.Passes = append(currentPackage.Passes, currentResult)
				lines = make([]string, 0, 0)
			} else if testName, ok := testFail(line); ok {
				currentResult = &Result{Name: testName}
				currentPackage.Fails = append(currentPackage.Fails, currentResult)
				lines = make([]string, 0, 0)
			} else if pkgName, ok := pgkDone(line); ok {
				log.Infof("Package done %s", pkgName)
				currentPackage.Name = pkgName
				results <- *currentPackage
				currentPackage = new(Package)
				lines = make([]string, 0, 0)
			} else {
				lines = append(lines, line)
				currentResult.Out = lines
			}
		}
		done <- struct{}{}
	}()
	return results, done
}

func testStarted(line string) (string, bool) {
	return search(line, "=== RUN")
}

func pgkDone(line string) (string, bool) {
	tName, tOK := search(line, "ok  \t")
	if tOK {
		return tName, tOK
	}
	return search(line, "FAIL\t")
}

func testSuccess(line string) (string, bool) {
	return search(line, "--- PASS:")
}

func testFail(line string) (string, bool) {
	return search(line, "--- FAIL:")
}

func search(line, search string) (string, bool) {
	index := strings.Index(line, search)
	if index < 0 {
		return "", false
	}
	stripped := line[index+len(search):]
	trimmed := strings.TrimSpace(stripped)
	tokens := strings.Fields(trimmed)
	return tokens[0], true
}

func runTests() (io.ReadCloser, chan struct{}) {
	done := make(chan struct{})
	pr, pw, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	go func() {
		defer close(done)
		defer func(done chan struct{}) { done <- struct{}{} }(done)
		defer pw.Close()

		ctx := context.Background()
		args := []string{"test", "./...", "-v"}
		args = append(args, os.Args[1:]...)
		c := exec.CommandContext(ctx, "go", args...)
		c.Stdout = pw
		c.Stderr = pw
		err := c.Run()
		if err != nil {
			log.Errorf("go test failed: %v", err)
			return
		}
		log.Infof("Idle")

	}()
	return pr, done
}
