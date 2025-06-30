package component

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var texts = []string{constant.Cancel, constant.Update, constant.Delete}

type Deal struct {
}

func NewDeal() *Deal {
	return &Deal{}
}

func (d *Deal) Content() fyne.CanvasObject {
	operators := make(map[int]*fyne.Container)
	var list *widget.List
	list = widget.NewList(
		func() int {
			return len(service.DataService.Deals)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel(""),
				container.NewHBox(widget.NewLabel(""), layout.NewSpacer(), widget.NewLabel("")),
				widget.NewLabel(""),
			)

		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			vbox := object.(*fyne.Container)
			deal := service.DataService.Deals[id]
			from, to := deal.AccountA, deal.AccountB
			cost := deal.AccountBPay
			if deal.AccountAPay > 0 {
				from, to = to, from
				cost = deal.AccountAPay
			}
			spent := strconv.FormatFloat(cost, 'f', 2, 64)
			label := widget.NewLabel(spent)
			label.Alignment = fyne.TextAlignTrailing
			date := deal.Date.Format("2006-01-02")
			objects := make([]fyne.CanvasObject, len(texts))
			var operator *fyne.Container
			for i, text := range texts {
				var f func()
				switch text {
				case constant.Cancel:
					f = func() {
						operator.Hide()
						list.Unselect(id)
					}
				case constant.Update:
					f = func() {
						updateItem <- deal
						service.DataService.RemoveDeal(deal)
						list.UnselectAll()
					}
				case constant.Delete:
					f = func() {
						service.DataService.RemoveDeal(deal)
						list.UnselectAll()
					}
				}
				objects[i] = widget.NewButton(text, f)
			}
			objects = append([]fyne.CanvasObject{layout.NewSpacer()}, objects...)
			operator = container.NewHBox(objects...)
			operator.Hide()
			operators[id] = operator

			vbox.Objects = []fyne.CanvasObject{
				widget.NewLabel(deal.Usage),
				container.NewHBox(widget.NewLabel(deal.Payee), layout.NewSpacer(), label),
				container.NewHBox(widget.NewLabel(fmt.Sprintf("%s %s-->%s", date, from, to)), layout.NewSpacer(), operator),
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		operators[id].Show()
	}
	return list
}
