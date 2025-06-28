package bill

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"log"
)

type Filter struct {
	content *widget.Select
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Content() fyne.CanvasObject {
	if f.content != nil {
		return f.content
	}
	accounts, err := service.DataService.Accounts.Get()
	if err != nil {
		log.Println(err)
	}
	accounts = append(accounts, constant.AccountPrefixes...)
	prefix := widget.NewSelect(accounts, func(account string) {
		service.DataService.FilterDataByAccount(account)
		//f.data.GetDeals(s)
		//f.data.RefreshPage()
	})
	prefix.SetSelected("")

	f.content = prefix
	return f.content
}
