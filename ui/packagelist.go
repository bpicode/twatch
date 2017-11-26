package ui

import (
	"fmt"
	"sort"
	"sync"

	"github.com/bpicode/twatch/test"
	ui "github.com/gizak/termui"
)

type PackageList struct {
	sync.Mutex
	*ui.List
	pkgMap map[string]test.Package
}

func (l *PackageList) Clear() {
	l.Lock()
	defer l.Unlock()
	l.pkgMap = make(map[string]test.Package)
	l.Items = nil
	l.Height = 2
}

func (l *PackageList) Update(p test.Package) {
	l.Lock()
	defer l.Unlock()
	l.pkgMap[p.Name] = p

	var pNames []string
	for pName := range l.pkgMap {
		pNames = append(pNames, pName)
	}
	sort.Strings(pNames)

	is := make([]string, 0, len(l.pkgMap))
	for _, pName := range pNames {
		is = append(is, l.fmtPkgResult(l.pkgMap[pName]))
	}
	l.Items = is
	l.Height = len(is) + 2
}

func (l *PackageList) fmtPkgResult(p test.Package) string {
	format := func(pk test.Package, style string) string {
		return fmt.Sprintf("[✔ %d ✘ %d %s](%s)", pk.NumPasses(), pk.NumFails(), pk.Name, style)
	}
	if !p.HasTests() {
		return format(p, "fg-white")
	}
	if p.HasErrors() {
		return format(p, "fg-red")
	}
	return format(p, "fg-green")
}

func NewPackageList() *PackageList {
	ls := &PackageList{List: ui.NewList(), pkgMap: make(map[string]test.Package)}
	ls.BorderLabel = "Packages"
	ls.BorderLabelFg = ui.ColorWhite & ui.AttrBold
	return ls
}
