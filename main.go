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
	service.DataRun()
	service.StatisticRun()
	service.Boot()

	//a.Preferences().SetStringList(constant.Period, []string{
	//	"2024-11",
	//	"2024-12",
	//	"2025-01",
	//	"2025-02",
	//	"2025-03",
	//	"2025-04",
	//	"2025-05",
	//	"2025-06",
	//})

	w.Resize(fyne.NewSize(600, 500))
	w.SetContent(page.NewIndex().Content())
	w.ShowAndRun()
}
