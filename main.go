package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"os"
	"path"
)

var (
	l        *widgets.List
	selected = hashset.New()
)

func main() {

	inProject, err := CheckInProjectRoot()
	if !inProject {
		fmt.Println("error:", fmt.Errorf("try executing from your project's root directory"))
		return
	}

	storage, err := NewStorageEngine()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("error:", fmt.Errorf("please specify the library you want to install as `dapp-pm [github-user]/[github-repo]@[github-version]`"))
		return
	}
	depString := os.Args[1]
	//specificSol := os.Args[2] // TODO: let power users shortcut the selection process if they know the exact name or path of the sol

	dep, err := storage.GetDependency(depString)
	if err != nil {
		panic(err)
	}

	extractor, err := storage.ExtractPaths(dep)

	if err != nil {
		panic(err)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l = widgets.NewList()
	l.Title = fmt.Sprintf(" Import Solidity files from [%s : %s] | a: toggle selection, c: import selected, q: quit ", dep.name, dep.version)

	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.SelectedRowStyle = ui.NewStyle(ui.ColorGreen)
	l.WrapText = false

	Resize()
	ReList(extractor)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if e.Type == ui.ResizeEvent {
			Resize()
		}
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		case "a":
			current := l.SelectedRow
			theSol := extractor.GetSol(current)
			if selected.Contains(theSol) {
				selected.Remove(theSol)
			} else {
				selected.Add(theSol)
			}
			ReList(extractor)
		case "c":
			storage.Commit(extractor, selected)
			return
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ReList(extractor)
	}

}

func Resize() {
	maxX, maxY := ui.TerminalDimensions()
	l.SetRect(0, 0, maxX, maxY)
}

func ReList(e *Extractor) {
	l.Rows = e.GetFileList(selected)
	ui.Render(l)
}

func CheckInProjectRoot() (bool, error) {
	working, err := os.Getwd()
	if err != nil {
		return false, err
	}
	dappFile, err := os.Open(path.Join(working, "Makefile")) // lazy / half-assed way to check we are in the root of the project
	if err != nil {
		return false, err
	}
	err = dappFile.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}
