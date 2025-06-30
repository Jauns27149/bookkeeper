package component

import (
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Date TODO 更新为中文显示
type Date struct {
	start *widget.DateEntry
	end   *widget.DateEntry
	all   *widget.Button
}

func (d *Date) Content() fyne.CanvasObject {
	return container.NewGridWithColumns(3, d.start, d.end, d.all)
}

func NewDate() *Date {
	server := service.DataService
	end := widget.NewDateEntry()
	end.Bind(server.End)
	start := widget.NewDateEntry()
	start.Bind(server.Start)

	end.OnSubmitted = func(s string) {
		server.LoadData()
	}

	all := widget.NewButton("all", func() {
		service.DataService.AllData()
	})

	return &Date{
		start: start,
		end:   end,
		all:   all,
	}
}
