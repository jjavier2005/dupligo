package ui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func NewMainWindow(app fyne.App, onScan func(string) ([][]string, error), onTrash func([]string) (int, error)) fyne.Window {
	window := app.NewWindow("DupliGo")

	var selectedPath string
	var checkboxes []*widget.Check
	var filePaths []string

	title := widget.NewLabelWithStyle(
		"DupliGo - Duplicate File Finder",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	folderLabel := widget.NewLabel("No hay ninguna carpeta seleccionada.")
	folderLabel.Wrapping = fyne.TextWrapWord

	statusLabel := widget.NewLabel("Selecciona una carpeta para comenzar.")
	selectionLabel := widget.NewLabel("0 carpetas seleccionadas.")

	resultsArea := container.NewVBox(widget.NewLabel("No hay resultados todavía"))

	deleteButton := widget.NewButton("Eliminar seleccionado", nil)

	deleteButton.Disable()

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

		if selectedCount > 0 {
			deleteButton.Enable()
		} else {
			deleteButton.Disable()
		}
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

	scanButton := widget.NewButton("Escanear", nil)

	clearResults := func(message string) {
		checkboxes = nil
		filePaths = nil

		resultsArea.Objects = []fyne.CanvasObject{
			widget.NewLabel(message),
		}
		resultsArea.Refresh()

		updateSelectionLabel()
		selectAllButton.Disable()
		deselectAllButton.Disable()
	}

	selectButton := widget.NewButton("Elige una carpeta", func() {
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
						filePaths = append(filePaths, path)

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

	deleteButton.OnTapped = func() {
		var selected []string

		// Find selected files.
		for i, check := range checkboxes {
			if check.Checked {
				selected = append(selected, filePaths[i])
			}
		}

		if len(selected) == 0 {
			return
		}

		dialog.ShowConfirm(
			"Confirmar envío a la papelera",
			fmt.Sprintf(
				"Vas a enviar %d archivos a la papelera.\n\n"+
					"Revisa cuidadosamente los archivos seleccionados.\n"+
					"Aunque SHA-256 es extremadamente fiable, "+
					"pueden existir errores o archivos modificados.\n\n"+
					"Podrías estar seleccionando todas las copias.\n\n"+
					"¿Quieres continuar?",
				len(selected),
			),
			func(confirmed bool) {
				if !confirmed {
					return
				}

				deleteButton.Disable()
				scanButton.Disable()
				selectButton.Disable()
				selectAllButton.Disable()
				deselectAllButton.Disable()

				for _, check := range checkboxes {
					check.Disable()
				}

				go func() {
					moved, err := onTrash(selected)

					fyne.Do(func() {
						clearResults("Vuelve a escanear para actualizar los resultados.")
						scanButton.Enable()
						selectButton.Enable()

						if err != nil {
							statusLabel.SetText(fmt.Sprintf(
								"%d archivos movidos antes del error", moved,
							))
							dialog.ShowError(err, window)
							return
						}

						statusLabel.SetText("Operación completada")
						dialog.ShowInformation(
							"Completado",
							fmt.Sprintf("%d archivos enviados a la papelera.", moved),
							window,
						)
					})
				}()
			},
			window,
		)
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
