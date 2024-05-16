package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func ErrorMessageQuit(err error, window fyne.Window, app fyne.App) {
	infoDialog := dialog.NewInformation("Ошибка", err.Error(), window)
	infoDialog.Show()
	infoDialog.SetOnClosed(func() {
		app.Quit()
	})
}
func ErrorMessage(err error, window fyne.Window) {
	infoDialog := dialog.NewInformation("Ошибка", err.Error(), window)
	infoDialog.Show()
}
