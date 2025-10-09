package component

import (
	"bookkeeper/constant-old"
	"bookkeeper/service-old"
	"bookkeeper/ui-old"
	"bookkeeper/util"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Picker struct {
	popup   *widget.PopUp
	content *fyne.Container
	c       *fyne.Container
	top     *fyne.Container
	confirm *widget.Button
	period  *fyne.Container
}

func (p *Picker) Content() fyne.CanvasObject {
	return widget.NewLabel("")
}

func (p *Picker) Popup() {
	if p.popup != nil && p.popup.Hidden {
		p.popup.Show()
	}

	var c fyne.CanvasObject

	var monthButton, custom *widget.Button
	monthButton = widget.NewButton("月份选择", func() {
		monthButton.Importance = widget.HighImportance
		custom.Importance = widget.MediumImportance
		monthScroll := ui.NewPicker(util.Months()).Content()
		yearScroll := ui.NewPicker(util.Years(5)).Content()
		temp := container.NewGridWithColumns(4, layout.NewSpacer(), yearScroll, monthScroll, layout.NewSpacer())
		p.content.Remove(p.c)
		p.content.Add(temp)
		p.c = temp
		p.popup.Refresh()
	})

	custom = widget.NewButton("自定义时间", func() {
		custom.Importance = widget.HighImportance
		monthButton.Importance = widget.MediumImportance
		year := ui.NewPicker(util.Years(5)).Content()
		month := ui.NewPicker(util.Months()).Content()
		n := time.Now()
		day := ui.NewPicker(util.Days(n.Year(), int(n.Month()))).Content()
		temp := container.NewGridWithColumns(5, layout.NewSpacer(), year, month, day, layout.NewSpacer())
		p.content.Remove(p.c)
		p.content.Add(temp)
		p.c = temp
		p.popup.Refresh()
	})
	top := container.NewVBox(container.NewHBox(monthButton, custom))

	confirm := widget.NewButton("Confirm", func() {
		service_old.BillService.DataEvent <- constant_old.Load
		p.popup.Hide()
	})
	confirm.Importance = widget.HighImportance

	p.content = container.NewBorder(top, confirm, nil, nil, c)

	canvas := util.PrimaryCanvas()
	popup := widget.NewPopUp(p.content, canvas)
	size := canvas.Size()
	popup.Resize(fyne.NewSize(size.Width, size.Height*0.618))
	popup.Move(fyne.NewPos(0, size.Height*0.382))
	popup.Show()

	p.popup = popup
}

func NewPicker() service_old.Component {
	return &Picker{}
}
