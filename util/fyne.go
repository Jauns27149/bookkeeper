package util

import "fyne.io/fyne/v2"

func PrimaryWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

func PrimaryCanvas() fyne.Canvas {
	return fyne.CurrentApp().Driver().AllWindows()[0].Canvas()
}
