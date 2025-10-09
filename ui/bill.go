package ui

import (
	"bookkeeper/constant"
	"bookkeeper/event"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _bill = &bill{
	statement: createStatementList(),
	aggregate: _aggregate._content,
	condition: getconditionContent(),
}

type bill struct {
	aggregate fyne.CanvasObject
	condition fyne.CanvasObject
	statement *widget.List
	content   fyne.CanvasObject
}

func init() {
	_bill.content = container.NewBorder(
		container.NewVBox(_bill.aggregate, _bill.condition),
		nil, nil, nil, _bill.statement,)

	setContent(constant.Bill, func() fyne.CanvasObject {
		return _bill.content
	})

	event.SetEventFunc(constant.BillRefresh, func() {
		_bill.content.Refresh()
	})
}
