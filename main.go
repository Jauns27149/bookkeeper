package main

import (
	"bookkeeper/service"
	"bookkeeper/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	w := app.NewWithID("bookkeeper").NewWindow("bookkeeper")
	w.Resize(fyne.NewSize(600, 500))

	service.Run()
	ui.Run()

	w.SetContent(ui.Content())
	w.ShowAndRun()
	
}
