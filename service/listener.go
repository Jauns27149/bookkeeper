package service

import (
	"bookkeeper/model"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"log"
	"sort"
	"time"
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

		m := make(map[time.Time][]model.Deal)
		for _, v := range data {
			if _, ok := m[v.Date]; ok {
				m[v.Date] = append(m[v.Date], v)
			} else {
				m[v.Date] = []model.Deal{v}
			}
		}
		s := make([]model.Statement, 0, len(m))
		for k, v := range m {
			s = append(s, model.Statement{Date: k, Deals: v})
		}
		sort.Slice(s, func(i, j int) bool {
			return s[i].Date.After(s[j].Date)
		})
		DataService.Statements = s

		DataService.Deals = data
		DataService.count()
		uiRefresh()
		log.Println("account change listened")
	}))
}

func uiRefresh() {
	fyne.CurrentApp().Driver().AllWindows()[0].Content().Refresh()
}
