package bill

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Filter struct {
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Content() fyne.CanvasObject {
	//accounts, err := service.DataService.Accounts.Get()
	//if err != nil {
	//	log.Println(err)
	//}
	//accounts = append(accounts, constant.AccountPrefixes...)
	//prefix := widget.NewSelectWithData(accounts, service.DataService.AccountType)
	//prefix.SetSelected(constant.All)
	prefix := widget.NewSelect([]string{"*"}, func(s string) {})
	return prefix
}
