package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewMainWindow(app fyne.App) fyne.Window {
	window := app.NewWindow("DupliGo")

	message := widget.NewLabel("Welcome to DupliGo!")

	button := widget.NewButton("Test interface", func() {
		message.SetText("The interface is working!")
	})

	window.SetContent(container.NewVBox(
		message,
		button,
	))

	window.Resize(fyne.NewSize(600, 400))

	return window
}
