package bill

import (
	"bookkeeper/constant-old"
	"bookkeeper/page/component"
	"bookkeeper/service-old"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type Head struct {
	gather service_old.Component
	period fyne.CanvasObject
	filter service_old.Component

	content fyne.CanvasObject
}

func (h *Head) Content() fyne.CanvasObject {
	if h.content != nil {
		return h.content
	}

	gather := h.gather.Content()
	h.content = container.NewVBox(gather, container.NewGridWithColumns(2, h.period, widget.NewButton("确定", func() {
		service_old.BillService.DataEvent <- constant_old.Load
	})), h.filter.Content())
	return h.content
}

func NewHead() *Head {
	period, err := service_old.BillService.Period.Get()
	if err != nil {
		log.Panicln(err)
	}
	text := fmt.Sprintf("%s ~ %s", period[0].Format(time.DateOnly), period[1].Format(time.DateOnly))
	button := widget.NewButton(text, func() {
		picker := component.Picker{}
		picker.Popup()
	})
	service_old.PageEventFunc[constant_old.Date] = func() {
		period, err := service_old.BillService.Period.Get()
		if err != nil {
			log.Panicln(err)
		}
		text := fmt.Sprintf("%s ~ %s", period[0].Format(time.DateOnly), period[1].Format(time.DateOnly))
		button.SetText(text)
		button.Refresh()

		log.Println("select date show update: ", text)
	}

	button.Alignment = widget.ButtonAlignLeading
	button.Importance = widget.HighImportance

	return &Head{
		period: container.NewHBox(button, layout.NewSpacer()),
		gather: component.NewGather(),
		filter: component.NewFilter(),
	}
}
