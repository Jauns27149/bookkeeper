package component

import (
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

// Date TODO 更新为中文显示
type Date struct {
	start *widget.DateEntry
	end   *widget.DateEntry
}

func (d *Date) Content() fyne.CanvasObject {
	return container.NewGridWithColumns(2, d.start, d.end)
}

func NewDate() *Date {
	now := time.Now()
	startTime := now.AddDate(0, 0, -now.Day()+1)
	end := widget.NewDateEntry()
	end.SetDate(&now)
	start := widget.NewDateEntry()
	start.SetDate(&startTime)

	end.OnSubmitted = func(s string) {
		ss, e := start.Date, end.Date
		service.DataService.ChangeDataByPeriod(*ss, *e)
	}

	return &Date{
		start: start,
		end:   end,
	}
}
