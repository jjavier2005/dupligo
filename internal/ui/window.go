package ui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewMainWindow(app fyne.App, onScan func(string) ([][]string, error)) fyne.Window {
	window := app.NewWindow("DupliGo")

	var selectedPath string
	var checkboxes []*widget.Check

	title := widget.NewLabelWithStyle(
		"DupliGo - Duplicate File Finder",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	folderLabel := widget.NewLabel("No hay ninguna carpeta seleccionada.")
	folderLabel.Wrapping = fyne.TextWrapWord

	statusLabel := widget.NewLabel("Selecciona una carpeta para comenzar.")
	selectionLabel := widget.NewLabel("0 carpetas seleccionadas.")

	resultsArea := container.NewVBox(
		widget.NewLabel("No hay resultados todavía"))

	updateSelectionLabel := func() {
		selectedCount := 0

		for _, check := range checkboxes {
			if check.Checked {
				selectedCount++
			}
		}

		selectionLabel.SetText(
			fmt.Sprintf("%d files selected", selectedCount),
		)
	}

	selectAllButton := widget.NewButton("Seleccionar todos", func() {
		for _, check := range checkboxes {
			check.SetChecked(true)
		}
		updateSelectionLabel()
	})

	deselectAllButton := widget.NewButton("Deseleccionar todos", func() {
		for _, check := range checkboxes {
			check.SetChecked(false)
		}
		updateSelectionLabel()
	})

	selectAllButton.Disable()
	deselectAllButton.Disable()

	deleteButton := widget.NewButton("Eliminar seleccionado", nil)

	scanButton := widget.NewButton("Escanear", nil)

	clearResults := func(message string) {
		checkboxes = nil

		resultsArea.Objects = []fyne.CanvasObject{
			widget.NewLabel(message),
		}
		resultsArea.Refresh()

		updateSelectionLabel()
		selectAllButton.Disable()
		deselectAllButton.Disable()
	}

	selectButton := widget.NewButton("Choose folder", func() {
		dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if folder == nil {
				return
			}

			selectedPath = folder.Path()
			folderLabel.SetText(selectedPath)

			statusLabel.SetText("Listo para escanear!")
			clearResults("Aún no hay resultados disponibles.")

			scanButton.Enable()
		}, window)
	})

	scanButton.OnTapped = func() {
		if selectedPath == "" {
			return
		}

		scanButton.Disable()
		selectButton.Disable()

		statusLabel.SetText("Escaneando archivos...")
		clearResults("Porfavor espere...")

		go func(path string) {
			groups, err := onScan(path)

			fyne.Do(func() {
				scanButton.Enable()
				selectButton.Enable()

				if err != nil {
					statusLabel.SetText("Escaneo fallido.")
					clearResults("No hay resultados disponibles.")
					dialog.ShowError(err, window)
					return
				}

				if len(groups) == 0 {
					statusLabel.SetText("No hay archivos coincidentes.")
					clearResults("No se encontraron duplicados.")
					return
				}

				accordion := widget.NewAccordion()
				totalFiles := 0

				for i, files := range groups {
					totalFiles += len(files)

					groupContent := container.NewVBox()

					for _, path := range files {
						check := widget.NewCheck(
							filepath.Base(path),
							func(bool) {
								updateSelectionLabel()
							},
						)

						checkboxes = append(checkboxes, check)

						pathLabel := widget.NewLabel(path)
						pathLabel.Wrapping = fyne.TextWrapWord

						fileRow := container.NewVBox(
							check,
							pathLabel,
							widget.NewSeparator(),
						)

						groupContent.Add(fileRow)
					}

					groupTitle := fmt.Sprintf(
						"Group %d (%d files)",
						i+1,
						len(files),
					)

					accordion.Append(
						widget.NewAccordionItem(
							groupTitle,
							groupContent,
						),
					)
				}

				// Show the first group initially.
				accordion.Open(0)

				resultsArea.Objects = []fyne.CanvasObject{
					accordion,
				}
				resultsArea.Refresh()

				statusLabel.SetText(fmt.Sprintf(
					"Found %d matching groups (%d files)",
					len(groups),
					totalFiles,
				))

				selectAllButton.Enable()
				deselectAllButton.Enable()
			})
		}(selectedPath)
	}

	header := container.NewVBox(
		title,
		folderLabel,
		container.NewHBox(selectButton, scanButton),
		statusLabel,
	)

	footer := container.NewBorder(
		nil,
		nil,
		container.NewHBox(selectAllButton, deselectAllButton),
		deleteButton,
		selectionLabel,
	)

	content := container.NewBorder(
		header,
		footer,
		nil,
		nil,
		container.NewVScroll(resultsArea),
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(900, 650))

	return window
}
