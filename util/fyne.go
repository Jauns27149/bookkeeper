package util

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"log"
)

func PrimaryWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

func PrimaryCanvas() fyne.Canvas {
	return fyne.CurrentApp().Driver().AllWindows()[0].Canvas()
}

func GetPrefString(pref binding.String)string{
	value, err := pref.Get()
	if err!=nil {
		log.Panic("pref is wrong")
	}
	return value
}
