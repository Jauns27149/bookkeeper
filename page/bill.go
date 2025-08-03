package page

import (
	"bookkeeper/page/bill"
	"bookkeeper/page/component"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"log"
)

type Bill struct {
	head service.Component
	//record  service.Component
	deal    service.Component
	filter  service.Component
	content fyne.CanvasObject
}

func (b *Bill) Content() fyne.CanvasObject {
	log.Println("create bill content .... ", b)
	if b.content != nil {
		return b.content
	}

	head := b.head.Content()
	//record := b.record.Content()
	deal := b.deal.Content()

	log.Println("create  content start.... ")
	b.content = container.NewBorder(
		container.NewVBox(head, b.filter.Content()),
		nil, nil, nil, deal)

	log.Println("create page content .... ", b)
	return b.content
}

func NewBill() *Bill {
	return &Bill{
		head:   bill.NewHead(),
		deal:   component.NewDeal(),
		filter: bill.NewFilter(),
	}
}
