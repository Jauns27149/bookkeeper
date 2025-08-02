package bill

import (
	"bookkeeper/page/component"
	"bookkeeper/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type Head struct {
	period fyne.CanvasObject
	gather service.Component
	date   fyne.CanvasObject
	//account service.Component
	content fyne.CanvasObject
}

func (h *Head) Content() fyne.CanvasObject {
	if h.content != nil {
		return h.content
	}

	gather := h.gather.Content()
	//account := h.account.Content()

	h.content = container.NewVBox(gather, h.period)
	return h.content
}

func NewHead() *Head {
	period, err := service.DataService.Period.Get()
	if err != nil {
		log.Panicln(err)
	}
	text := fmt.Sprintf("%s ~ %s", period[0].Format(time.DateOnly), period[1].Format(time.DateOnly))
	button := widget.NewButton(text, func() {
		picker := component.Picker{}
		picker.Popup()
	})
	button.Alignment = widget.ButtonAlignLeading
	button.Importance = widget.HighImportance

	return &Head{
		period: container.NewHBox(button, layout.NewSpacer()),
		gather: component.NewGather(),
		//account: NewAccount(),
		//date: component.NewDate().Content(),
	}
}
