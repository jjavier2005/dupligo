package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	application := app.New()

	window := application.NewWindow("DupliGo")

	message := widget.NewLabel("Welcome to DupliGo!")

	button := widget.NewButton("Test interface", func() {
		message.SetText("The interface is working!")
	})

	window.SetContent(container.NewVBox(
		message,
		button,
	))

	window.Resize(fyne.NewSize(600, 400))

	window.ShowAndRun()
}
