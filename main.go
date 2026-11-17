package main

import (
	"sort"

	"fyne.io/fyne/v2/app"

	"github.com/jjavier2005/dupligo/internal/duplicates"
	"github.com/jjavier2005/dupligo/internal/recycle"
	"github.com/jjavier2005/dupligo/internal/scanner"
	"github.com/jjavier2005/dupligo/internal/ui"
)

func runScan(root string) ([][]string, error) {

	paths, err := scanner.ListFiles(root)
	if err != nil {
		return nil, err
	}

	sizeMap, err := duplicates.GetRepeatedSizes(paths)
	if err != nil {
		return nil, err
	}

	hashMap, err := duplicates.GetSameHashFiles(sizeMap)
	if err != nil {
		return nil, err
	}

	var duplicatesList [][]string

	for _, files := range hashMap {
		duplicatesList = append(duplicatesList, files)
	}

	sort.Slice(duplicatesList, func(i, j int) bool { return len(duplicatesList[i]) > len(duplicatesList[j]) })

	return duplicatesList, nil
}

func main() {
	application := app.NewWithID("com.jjavier2005.dupligo")
	window := ui.NewMainWindow(application, runScan, recycle.Send)
	window.ShowAndRun()
}
