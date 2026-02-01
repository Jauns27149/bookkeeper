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

type mapAccountsObject struct {
	name   *widget.Label
	amount *widget.Label
	option *HideOption
}

type HideOption struct {
	content *fyne.Container

	rename *widget.Button
}

func (a *account) createContent() {
	if a.content != nil {
		return
	}

	grid := container.NewGridWithColumns(len(a.buttons))
	for _, b := range a.buttons {
		grid.Add(b)
	}

	bottom := container.NewHBox(layout.NewSpacer(), a.add)
	a.content = container.NewBorder(grid, bottom, nil, nil, a.list)
	log.Println("create accounts content finished")
}

func (a *account) createList() {
	store := model.NewListStore()
	_accounts.list = widget.NewList(
		func() int { return len(source) },
		func() fyne.CanvasObject {
			acc := &mapAccountsObject{
				name:   widget.NewLabel(constant.Zero),
				amount: widget.NewLabel(constant.Zero),
				option: newHideOption(),
			}
			hbox := container.NewHBox(acc.name, layout.NewSpacer(), acc.amount)
			c := container.NewVBox(hbox, acc.option.content)
			acc.option.content.Hide()
			store.Store(c, acc)
			return c
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			store.Store(id, obj)
			acc := store.Load(obj).(*mapAccountsObject)
			acc.name.SetText(source[id].Name)
			acc.amount.SetText(strconv.FormatFloat(source[id].Amount, 'f', 2, 64))
		},
	)

	_accounts.list.OnSelected = func(id widget.ListItemID) {
		acc := store.Load(id).(*mapAccountsObject)
		acc.option.content.Show()
		_accounts.list.SetItemHeight(id, acc.name.Size().Height*2+5)

		acc.option.rename.OnTapped = func() {
			entry := widget.NewEntry()
			button := widget.NewButton("确定", nil)
			entry.PlaceHolder = "			"
			c := container.NewBorder(nil, nil, nil, button, entry)

			up := widget.NewPopUp(c, currentCanvas())
			size := up.Size()
			size.Width *= 3
			up.Resize(size)

			pre := currentCanvas().Size()
			curr := up.MinSize()
			up.Move(fyne.NewPos(pre.Width/2-curr.Width/2, pre.Width/2-curr.Width/2))
			up.Show()

			button.OnTapped = func() {
				service.GetAccounts().Rename(acc.name.Text, entry.Text)
				up.Hide()
			}
			_accounts.list.OnSelected(id)
		}
	}

	_accounts.list.OnUnselected = func(id widget.ListItemID) {
		acc := store.Load(id).(*mapAccountsObject)
		acc.option.content.Hide()
		_accounts.list.SetItemHeight(id, acc.name.Size().Height)
	}
}

func newHideOption() *HideOption {
	option := &HideOption{
		rename: widget.NewButton("账号重命名", nil),
	}
	option.rename.Importance = widget.HighImportance

	option.content = container.NewHBox(layout.NewSpacer(), option.rename)
	return option
}

func (a *account) createBottom() {
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
		a.createBottom()

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
