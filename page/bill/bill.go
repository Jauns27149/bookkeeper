package bill

import (
	"bookkeeper/intf"
	"bookkeeper/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var updateItem = make(chan model.Deal, 1)

type Bill struct {
	head    intf.Component
	record  intf.Component
	deal    intf.Component
	content fyne.CanvasObject
	filter  intf.Component
}

func NewBill() *Bill {
	return &Bill{
		head:   NewHead(),
		record: NewRecord(),
		deal:   NewDeal(),
		filter: NewFilter(),
	}
}

func (b *Bill) Content() fyne.CanvasObject {
	if b.content != nil {
		return b.content
	}

	head := b.head.Content()
	record := b.record.Content()
	deal := b.deal.Content()

	b.content = container.NewBorder(
		container.NewVBox(head, record, b.filter.Content()),
		nil, nil, nil, deal)
	return b.content
}
