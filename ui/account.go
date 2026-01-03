package ui

import (
	"bookkeeper/constant"
	"bookkeeper/model"
	"bookkeeper/service"
	"fmt"
	"log"
	"maps"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var _accounts = &account{}
var source = []model.AccountDetail{}

type account struct {
	content fyne.CanvasObject
	buttons []*widget.Button
	list    *widget.List
	detail  *widget.List
	add     *widget.Button
}

type mapAcountsObject struct {
	name   *widget.Label
	amount *widget.Label
}

func (a *account) createContent() {
	if a.content != nil {
		return
	}

	grid := container.NewGridWithColumns(len(a.buttons))
	for _, b := range a.buttons {
		grid.Add(b)
	}

	buttom := container.NewHBox(layout.NewSpacer(), a.add)
	a.content = container.NewBorder(grid, buttom, nil, nil, a.list)
	log.Println("create accounts content finished")
}

func (a *account) createList() {
	m := make(map[fyne.CanvasObject]*mapAcountsObject)
	_accounts.list = widget.NewList(
		func() int { return len(source) },
		func() fyne.CanvasObject {
			account := &mapAcountsObject{widget.NewLabel(constant.Zero), widget.NewLabel(constant.Zero)}
			c := container.NewHBox(account.name, layout.NewSpacer(), account.amount)
			m[c] = account
			return c
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			m[object].name.SetText(source[id].Name)
			m[object].amount.SetText(strconv.FormatFloat(source[id].Amount, 'f', 2, 64))
		},
	)
}

func (a *account) createButtom() {
	a.add = widget.NewButton(constant.AddAccount, func() {
		var popup *widget.PopUp

		entry := widget.NewEntry()
		strs := slices.Sorted(maps.Keys(service.GetCondition().Account))
		prefix := widget.NewSelect(strs, nil)
		sure := widget.NewButton(constant.Sure, func() {
			account := fmt.Sprintf("%s:%s", prefix.Selected, entry.Text)
			service.AddAccount(account)
			popup.Hide()
		})

		c := container.NewHBox(prefix, entry, sure)

		popup = widget.NewPopUp(c, currentCanvas())
		pre := currentCanvas().Size()
		curr := popup.MinSize()
		popup.Move(fyne.NewPos(pre.Width/2-curr.Width/2, pre.Width/2-curr.Width/2))
		popup.Show()
	})

}

func (a *account) run() {
	go func() {
		flag := make(chan struct{})
		a.createList()
		a.createButtom()

		data := service.GetAccounts().Accounts
		for _, item := range data {
			a.buttons = append(a.buttons, widget.NewButton(item.Category, func() {
				source = item.AccountDetail
				a.list.Refresh()
			}))
		}
		if len(data) >= 4 {
			source = data[4].AccountDetail
		}

		close(flag)
		setContent(constant.Account, func() fyne.CanvasObject {
			<-flag
			a.createContent()
			return a.content
		})
	}()
}
