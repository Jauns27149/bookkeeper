package service

import (
	"bookkeeper/model"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"log"
)

func preference() {
	fyne.CurrentApp().Preferences().AddChangeListener(func() {
		DataService.LoadData()

		uiRefresh()
	})
}

func account() {
	accountType := DataService.AccountType
	accountType.AddListener(binding.NewDataListener(func() {
		acc, err := accountType.Get()
		if err != nil {
			log.Println(err)
			return
		}
		data := make([]model.Deal, 0)
		for _, deal := range DataService.Deals {
			if util.CheckAccount(deal, acc) {
				data = append(data, deal)
			}
		}
		DataService.Deals = data
		DataService.count()
		uiRefresh()
		log.Println("account change listened")
	}))
}

func uiRefresh() {
	fyne.CurrentApp().Driver().AllWindows()[0].Content().Refresh()
}
