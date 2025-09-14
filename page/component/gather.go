package component

import (
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

type Gather struct {
}

func NewGather() *Gather {
	return &Gather{}
}

func (g *Gather) Content() fyne.CanvasObject {
	log.Println("gather content start")
	log.Println("get server: ")
	head := []string{"收入", "支出", "负债", "预算"}
	grib := container.NewGridWithColumns(len(head))
	for _, h := range head {
		log.Println(h)
		grib.Add(widget.NewLabel(h))
	}
	grib.Add(widget.NewLabelWithData(service.BillService.Income))
	grib.Add(widget.NewLabelWithData(service.BillService.Expense))
	grib.Add(widget.NewLabelWithData(service.BillService.Liability))
	grib.Add(widget.NewLabelWithData(service.BillService.Budget))

	for i, o := range grib.Objects {
		label := o.(*widget.Label)
		label.Alignment = fyne.TextAlignCenter
		if i >= len(head) {
			label.Importance = widget.WarningImportance
		}
	}
	log.Println("gather content successful")
	return grib
}
