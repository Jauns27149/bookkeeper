package component

import (
    "bookkeeper/constant"
    "bookkeeper/layoutCustom"
    "bookkeeper/service"
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
    preList := []string{constant.Expenses, constant.Liabilities, constant.Assets, constant.Income, constant.All}
    prefix := widget.NewSelectWithData(preList, service.BillService.Prefix)
    suffix := widget.NewSelectWithData([]string{}, service.BillService.Suffix)
    return &Filter{Prefix: prefix, Suffix: suffix}
}

func (f *Filter) Content() fyne.CanvasObject {
    service.BillService.Prefix.AddListener(binding.NewDataListener(func() {
        prefix, _ := service.BillService.Prefix.Get()
        f.Suffix.SetOptions(service.AccountService.AccountMap()[prefix])
        f.Suffix.SetSelectedIndex(0)
    }))
    f.Prefix.SetSelected(constant.All)

    return container.New(layoutCustom.NewSplit(0.3), f.Prefix, f.Suffix)
}
