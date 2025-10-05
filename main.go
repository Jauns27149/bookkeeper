package main

import (
	"bookkeeper/event"
	"bookkeeper/page"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	log.Println("账本开始启动...")
	a := app.NewWithID("bookkeeper")
	w := a.NewWindow("bookkeeper")

	service.Boot()
	event.Run()
	log.Println("data start successful")

	w.Resize(fyne.NewSize(600, 500))
	w.SetContent(page.NewIndex().Content())
	w.ShowAndRun()

}
