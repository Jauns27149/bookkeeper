package component

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"bookkeeper/util"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"time"
)

type Deal struct {
}

func NewDeal() *Deal {
	return &Deal{}
}

func (d *Deal) Content() fyne.CanvasObject {
	var list *widget.List
	var currentBlabel *widget.Label
	var currentStack *fyne.Container
	list = widget.NewList(
		func() int {
			return len(service.DataService.Statements)
		},
		func() fyne.CanvasObject {
			return container.NewVBox()
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			statement := service.DataService.Statements[id]
			if id == 5 {
				fmt.Println(service.DataService.Statements)
			}
			title := widget.NewLabel(statement.Date.Format(time.DateOnly))
			vbox := container.NewVBox(title)

			for _, item := range statement.Deals {
				from, to := item.Payment.Name, item.Receiver.Name
				cost := item.Receiver.Cost
				if item.Payment.Cost > 0 {
					from, to = to, from
					cost = item.Payment.Cost
				}
				richText := widget.NewRichText(
					&widget.TextSegment{
						Text: fmt.Sprintf("%s %s", item.Payee, item.Usage),
						Style: widget.RichTextStyle{
							ColorName: theme.ColorNameSuccess,
						},
					},
					&widget.TextSegment{
						Text: fmt.Sprintf("%s >> %s", util.LastAccount(from), util.LastAccount(to)),
						Style: widget.RichTextStyle{
							ColorName: theme.ColorNameSuccess,
						},
					},
				)
				costText := widget.NewLabel(strconv.FormatFloat(cost, 'f', 2, 64))
				subStatement := container.NewHBox(richText, layout.NewSpacer(), costText)

				buttonsH := container.NewHBox()
				update := widget.NewButton(constant.Update, func() {
					service.TallyService.Finish <- constant.Tally
					service.DataService.RemoveDeal(item)
					buttonsH.Hide()
				})
				update.Importance = widget.WarningImportance

				deleteItem := widget.NewButton(constant.Delete, func() {
					service.DataService.RemoveDeal(item)
					buttonsH.Hide()
				})
				deleteItem.Importance = widget.DangerImportance
				buttonsH.Add(layout.NewSpacer())
				buttonsH.Add(update)
				buttonsH.Add(deleteItem)

				var stack *fyne.Container
				var button *widget.Button
				button = widget.NewButton("", func() {
					if currentBlabel != nil {
						currentBlabel.Show()
						currentStack.Objects[0], currentStack.Objects[1] = currentStack.Objects[1], currentStack.Objects[0]
					}
					if currentBlabel != costText {
						costText.Hide()
						stack.Objects[0], stack.Objects[1] = stack.Objects[1], stack.Objects[0]
						currentBlabel = costText
						currentStack = stack
					} else {
						currentBlabel = nil
						currentStack = nil
					}
					log.Println("operation finished")
				})
				stack = container.NewStack(buttonsH, button, subStatement)
				vbox.Add(stack)
			}
			c := object.(*fyne.Container)
			c.RemoveAll()
			c.Add(vbox)

			height := vbox.MinSize().Height
			list.SetItemHeight(id, height)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		list.UnselectAll()
	}

	return list
}
