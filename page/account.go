package page

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
)

type Account struct {
	subContents    []fyne.CanvasObject
	categorization []*widget.Button
	current        int
}

func (a *Account) Content() fyne.CanvasObject {
	var content *fyne.Container
	for i, button := range a.categorization {
		ii := i
		button.OnTapped = func() {
			content.Remove(a.subContents[a.current])
			content.Add(a.subContents[ii])
			a.categorization[a.current].Enable()
			a.categorization[ii].Disable()
			a.current = ii
		}
	}

	objects := make([]fyne.CanvasObject, len(a.subContents))
	for i, button := range a.categorization {
		objects[i] = button
	}
	top := container.NewGridWithColumns(len(a.categorization), objects...)

	a.current = 3
	a.categorization[a.current].Disable()
	content = container.NewBorder(top, nil, nil, nil, a.subContents[a.current])
	return content
}

func NewAccount() *Account {
	prefixes := []string{constant.Expenses, constant.Income, constant.Assets, constant.Liabilities}
	categorization, components := make([]*widget.Button, len(prefixes)), make([]fyne.CanvasObject, len(prefixes))
	accounts := service.AccountService.AccountCollection()
	accountsMap := util.GroupAccountByPrefix(accounts)

	for i, prefix := range prefixes {
		categorization[i] = widget.NewButton(prefix, nil)
		subAccounts := accountsMap[prefix]
		rows := make([]fyne.CanvasObject, len(subAccounts))
		for ii, account := range subAccounts {
			head := canvas.NewText(account.Account, color.White)
			tail := canvas.NewText(strconv.FormatFloat(account.Amount, 'f', 2, 64), color.White)
			rows[ii] = container.NewHBox(head, layout.NewSpacer(), tail)
		}
		components[i] = container.NewVBox(rows...)
	}

	return &Account{
		subContents:    components,
		categorization: categorization,
	}
}
