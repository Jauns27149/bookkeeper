package ui

import (
	"bookkeeper/layoutCustom"
	"bookkeeper/model"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
)

func AccountSelect(account model.BindAccount, pre string) fyne.CanvasObject {
	prefixes := service.TallyService.Prefixes()
	head := widget.NewSelectWithData(prefixes, account.Prefix)
	tail := widget.NewSelectWithData(nil, account.Account)
	err := account.Prefix.Set(pre)
	if err != nil {
		log.Panicln(err)
	}
	account.Prefix.AddListener(binding.NewDataListener(func() {
		prefix, err := account.Prefix.Get()
		if err != nil {
			log.Panicln(err)
		}
		tail.SetOptions(service.TallyService.Suffixes(prefix))
		tail.SetSelectedIndex(0)
		tail.Refresh()
	}))

	return container.New(layoutCustom.NewSplit(0.318), head, tail)
}
