package bill

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"log"
)

type Filter struct {
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Content() fyne.CanvasObject {
	accounts, err := service.DataService.Accounts.Get()
	if err != nil {
		log.Println(err)
	}
	accounts = append(accounts, constant.AccountPrefixes...)
	prefix := widget.NewSelectWithData(accounts, service.DataService.AccountType)
	prefix.SetSelected(constant.All)

	return prefix
}
