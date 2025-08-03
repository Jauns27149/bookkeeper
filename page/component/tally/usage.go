package tally

import (
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type usage struct {
	payee   *widget.Entry
	utilize *widget.Entry
}

func (u *usage) Content() fyne.CanvasObject {
	return container.NewGridWithColumns(2, u.payee, u.utilize)
}

func NewUsage() service.Component {
	data := service.BillService
	payee := widget.NewEntryWithData(data.Payee)
	utilize := widget.NewEntryWithData(data.Utilize)

	return &usage{
		payee:   payee,
		utilize: utilize,
	}
}
