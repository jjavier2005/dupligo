package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewMainWindow(app fyne.App) fyne.Window {
	window := app.NewWindow("DupliGo")

	title := widget.NewLabel("DupliGo")

	message := widget.NewLabel("Welcome to DupliGo!")

	selectButton := widget.NewButton("Choose folder", func() {
		dialog.ShowFolderOpen(
			func(folder fyne.ListableURI, err error) {
				if err != nil {
					dialog.ShowError(err, window)
					return
				}

				if folder != nil {
					message.SetText("Selected folder: " + folder.Path())
				}
			},

			window,
		)
	})

	window.SetContent(container.NewVBox(
		title,
		message,
		selectButton,
	))

	window.Resize(fyne.NewSize(600, 400))

	return window
}
