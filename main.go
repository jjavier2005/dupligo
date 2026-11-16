package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/jjavier2005/dupligo/internal/ui"
)

func main() {
	application := app.New()
	window := ui.NewMainWindow(application)
	window.ShowAndRun()
}
