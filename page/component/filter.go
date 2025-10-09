package component

import (
	"bookkeeper/constant-old"
	"bookkeeper/layoutCustom"
	"bookkeeper/service-old"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Filter struct {
	Prefix *widget.Select
	Suffix *widget.Select
}

func NewFilter() *Filter {
	preList := []string{constant_old.Expenses, constant_old.Liabilities, constant_old.Assets, constant_old.Income, constant_old.All}
	prefix := widget.NewSelectWithData(preList, service_old.BillService.Prefix)
	suffix := widget.NewSelectWithData([]string{}, service_old.BillService.Suffix)
	return &Filter{Prefix: prefix, Suffix: suffix}
}

func (f *Filter) Content() fyne.CanvasObject {
	service_old.BillService.Prefix.AddListener(binding.NewDataListener(func() {
		prefix, _ := service_old.BillService.Prefix.Get()
		f.Suffix.SetOptions(service_old.AccountService.AccountMap()[prefix])
		f.Suffix.SetSelectedIndex(0)
	}))
	f.Prefix.SetSelected(constant_old.All)

	return container.New(layoutCustom.NewSplit(0.3), f.Prefix, f.Suffix)
}
