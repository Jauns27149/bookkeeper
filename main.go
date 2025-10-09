package main

import (
	"bookkeeper/app"

	"bookkeeper/service"
	"bookkeeper/ui"
)

func main() {
	app.Run()
	service.Run()
	ui.Run()
	content := ui.Content()
	app.Window.SetContent(content)
	app.Window.ShowAndRun()

	// w :=app.New().NewWindow("bookkeeper")
	// w.SetContent(widget.NewLabel("你好"))
	// w.ShowAndRun()
}
