package ui

import "fyne.io/fyne/v2"

func Run() {
	_home.run()
}

func currentCanvas() fyne.Canvas {
	return fyne.CurrentApp().Driver().AllWindows()[0].Canvas()
}
