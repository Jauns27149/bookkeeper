package bill

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Account struct {
	content fyne.CanvasObject
	data    *service.Data
	button  *widget.Button
	entry   *widget.Entry
	sel     *widget.Select
}

func NewAccount() *Account {
	data := service.DataService
	button := widget.NewButton(constant.CreateAccount, nil)
	entry := widget.NewEntry()
	sel := widget.NewSelect(constant.AccountPrefixes, nil)

	sel.SetSelected(constant.Expenses)
	entry.SetPlaceHolder(constant.InputAccountName)

	return &Account{
		data:   data,
		button: button,
		entry:  entry,
		sel:    sel,
	}
}

func (a *Account) Content() fyne.CanvasObject {
	if a.content != nil {
		return a.content
	}

	grid := container.NewGridWithColumns(2, a.sel, a.entry)
	grid.Hide()
	a.button.OnTapped = func() {
		if grid.Hidden {
			grid.Show()
		} else {
			grid.Hide()
		}
	}
	a.entry.OnSubmitted = func(s string) {
		pref := a.data.Pref
		accounts := pref.StringList(constant.Accounts)
		value := fmt.Sprintf("%v:%v", a.sel.Selected, a.entry.Text)
		accounts = append(accounts, value)
		pref.SetStringList(constant.Accounts, accounts)
		grid.Hide()
		a.data.Refresh()
	}

	a.content = container.NewVBox(a.button, grid)
	return a.content
}
