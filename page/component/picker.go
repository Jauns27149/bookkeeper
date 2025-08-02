package component

import (
	"bookkeeper/service"
	"bookkeeper/ui"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Picker struct{}

func (p *Picker) Content() fyne.CanvasObject {
	return widget.NewLabel("")
}

func (p *Picker) Popup() {
	monthButton := widget.NewButton("Select Month", nil)
	monthButton.OnTapped = func() {
		monthButton.Importance = widget.HighImportance
	}
	top := container.NewHBox(monthButton)

	monthScroll := ui.NewPicker(util.Months(3)).Content()
	yearScroll := ui.NewPicker(util.Years(5)).Content()

	confirm := widget.NewButton("Confirm", func() {})
	confirm.Importance = widget.HighImportance

	content := container.NewBorder(top, confirm, nil, nil,
		container.NewGridWithColumns(4, layout.NewSpacer(), yearScroll, monthScroll, layout.NewSpacer()),
	)

	canvas := util.PrimaryCanvas()
	popup := widget.NewPopUp(content, canvas)
	size := canvas.Size()
	popup.Resize(fyne.NewSize(size.Width, size.Height*0.618))
	popup.Move(fyne.NewPos(0, size.Height*0.382))
	popup.Show()
}

func NewPicker() service.Component {
	return &Picker{}
}
