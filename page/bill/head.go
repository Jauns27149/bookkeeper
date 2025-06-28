package bill

import (
	"bookkeeper/intf"
	"bookkeeper/page/component"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Head struct {
	gather  intf.Component
	date    fyne.CanvasObject
	account intf.Component
	content fyne.CanvasObject
}

func NewHead() *Head {
	return &Head{
		gather:  NewGather(),
		account: NewAccount(),
		date:    component.NewDate().Content(),
	}
}

func (h *Head) Content() fyne.CanvasObject {
	if h.content != nil {
		return h.content
	}

	gather := h.gather.Content()
	account := h.account.Content()

	h.content = container.NewVBox(gather, h.date, account)
	return h.content
}
