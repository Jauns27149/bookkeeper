package bill

import (
	"bookkeeper/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

type Deal struct {
	content *widget.List
	data    *service.Data
}

func NewDeal() *Deal {
	return &Deal{
		data: service.DataService,
	}
}

// Content TODO need to split
// TODO list组件无法取消被选择状态
func (d *Deal) Content() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	buttons := make([]*widget.Button, len(d.data.Deals))
	updates := make([]*widget.Button, len(d.data.Deals))
	d.content = widget.NewList(
		func() int {
			if len(buttons) < len(d.data.Deals) {
				buttons = append(buttons, nil)
				updates = append(updates, nil)
			}
			return len(d.data.Deals)
		},

		func() fyne.CanvasObject {
			//return widget.NewLabel("")
			items := make([]fyne.CanvasObject, 3)
			for i := range 3 {
				items[i] = canvas.NewText("", nil)
			}
			h := container.NewVBox(items...)

			items = make([]fyne.CanvasObject, 5)
			items[0] = h
			items[1] = layout.NewSpacer()
			items[2] = canvas.NewText("", nil)
			items[3] = widget.NewButton("更新", nil)
			items[4] = widget.NewButton("删除", nil)
			v := container.NewHBox(items...)
			return v
		},

		func(id widget.ListItemID, object fyne.CanvasObject) {
			v := object.(*fyne.Container)
			left := v.Objects[0].(*fyne.Container)
			pay := v.Objects[2].(*canvas.Text)
			update := v.Objects[3].(*widget.Button)
			button := v.Objects[4].(*widget.Button)

			deal := d.data.Deals[id]
			left.Objects[0].(*canvas.Text).Text = deal.Usage
			left.Objects[1].(*canvas.Text).Text = deal.Payee
			head, tail := "", ""
			if deal.AccountAPay < 0 {
				head, tail = deal.AccountA, deal.AccountB
			} else {
				head, tail = deal.AccountB, deal.AccountA
			}
			text := fmt.Sprintf("%v %v->%v", deal.Date.Format(time.DateOnly), head, tail)
			left.Objects[2].(*canvas.Text).Text = text
			pay.Text = strconv.FormatFloat(deal.AccountAPay, 'f', 2, 64)
			button.Importance = widget.WarningImportance
			button.Hide()
			update.Hide()
			button.OnTapped = func() {
				d.data.RemoveDeal(d.data.Deals[id])
			}
			update.OnTapped = func() {
				updateItem <- d.data.Deals[id]
				button.OnTapped()
			}

			buttons[id] = button
			updates[id] = update
		},
	)

	d.content.OnSelected = func(id widget.ListItemID) {
		buttons[id].Show()
		updates[id].Show()
	}
	return d.content
}
