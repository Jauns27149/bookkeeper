package component

import (
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Gather struct {
}

func NewGather() *Gather {
	return &Gather{}
}

func (g *Gather) Content() fyne.CanvasObject {
	server := service.DataService
	head := []string{"收入", "支出", "负债"}
	grib := container.NewGridWithColumns(len(head))
	for _, h := range head {
		grib.Add(widget.NewLabel(h))
	}
	grib.Add(widget.NewLabelWithData(server.Income))
	grib.Add(widget.NewLabelWithData(server.Expense))
	grib.Add(widget.NewLabelWithData(server.Liability))

	for _, o := range grib.Objects {
		o.(*widget.Label).Alignment = fyne.TextAlignCenter
	}

	return grib
}
