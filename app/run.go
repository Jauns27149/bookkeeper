package app

import (
	"fyne.io/fyne/v2/theme"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var Window fyne.Window

var flagAppRun = make(chan struct{})

func Preferences() fyne.Preferences {
	<-flagAppRun
	return fyne.CurrentApp().Preferences()
}

func Run() {
	// theme.DefaultTheme().Font()
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
}

func init() {
	Window = app.NewWithID("bookkeeper").NewWindow("bookkeeper")
	Window.Resize(fyne.NewSize(600, 500))

	close(flagAppRun)
	log.Println("app start successful")
}
