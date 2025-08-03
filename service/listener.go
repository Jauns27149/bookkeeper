package service

import (
	"bookkeeper/constant"
	"bookkeeper/util"
)

func dataEvent() {
	for {
		switch <-BillService.DataEvent {
		case constant.Load:
			go BillService.Load()
		case constant.Count:
			statements := BillService.Statements
			m := util.CountData(statements)
			var budget float64
			for key, value := range m {
				switch key {
				case constant.Income:
					value = -value
					util.SetValue(BillService.Head.Income, value)
					budget += value * 0.618
				case constant.Expenses:
					value = -value
					util.SetValue(BillService.Head.Expense, value)
					budget += value
				case constant.Liabilities:
					util.SetValue(BillService.Head.Liability, value)
				}
			}
			util.SetValue(BillService.Head.Budget, budget)
			BillService.UiEvent <- constant.Index
		}
	}
}

func accountBind() {

}

//func preference() {
//	fyne.CurrentApp().Preferences().AddChangeListener(func() {
//		BillService.LoadData()
//
//		uiRefresh()
//	})
//}
//
//func account() {
//	accountType := BillService.AccountType
//	accountType.AddListener(binding.NewDataListener(func() {
//		acc, err := accountType.Get()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		data := make([]model.Deal, 0)
//		for _, deal := range BillService.Deals {
//			if util.CheckAccount(deal, acc) {
//				data = append(data, deal)
//			}
//		}
//
//		m := make(map[time.Time][]model.Deal)
//		for _, v := range data {
//			if _, ok := m[v.Date]; ok {
//				m[v.Date] = append(m[v.Date], v)
//			} else {
//				m[v.Date] = []model.Deal{v}
//			}
//		}
//		s := make([]model.Statement, 0, len(m))
//		for k, v := range m {
//			s = append(s, model.Statement{Date: k, Deals: v})
//		}
//		sort.Slice(s, func(i, j int) bool {
//			return s[i].Date.After(s[j].Date)
//		})
//		BillService.Statements = s
//
//		BillService.Deals = data
//		BillService.count()
//		uiRefresh()
//		log.Println("account change listened")
//	}))
//}
//
//func uiRefresh() {
//	fyne.CurrentApp().Driver().AllWindows()[0].Content().Refresh()
//}
