package ui

import (
	"bookkeeper/constant"
	"bookkeeper/event"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _bill = &bill{}

type bill struct {
	aggregate fyne.CanvasObject
	condition fyne.CanvasObject
	statement *widget.List
	content   fyne.CanvasObject
}

func (b *bill) run() {
	_aggregate.run()

	b.statement = createStatementList()
	b.aggregate = _aggregate.content
	b.condition = getconditionContent()
	b.content = container.NewBorder(
		container.NewVBox(_bill.aggregate, _bill.condition),
		nil, nil, nil, _bill.statement)

	setContent(constant.Bill, func() fyne.CanvasObject {
		return b.content
	})

	event.SetEventFunc(constant.BillRefresh, func() {
		b.content.Refresh()
	})


}
