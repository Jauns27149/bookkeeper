package main

import (
	"bookkeeper/page"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
)

func main() {
	log.Println("账本开始启动...")
	a := app.NewWithID("bookkeeper")
	w := a.NewWindow("bookkeeper")

	// 初始化数据
	service.Boot()
	service.StatisticRun()

	w.Resize(fyne.NewSize(600, 500))
	w.SetContent(page.NewIndex().Content())
	w.ShowAndRun()
}
