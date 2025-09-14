package service

import (
	"bookkeeper/constant"
	"bookkeeper/event"
	"bookkeeper/util"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"log"
	"time"
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
			if event.CurrentEvent != constant.UpdateEvent {
				BillService.UiEvent <- constant.Index
			}
		}
	}
}

var PageEventFunc = make(map[uint]func())

func uiEvent() {
	for {
		key := <-BillService.UiEvent
		fmt.Println("ui event: ", key)
		fyne.Do(func() {
			if fn, ok := PageEventFunc[uint(key)]; ok {
				fn()
			}
		})
	}
}

func AccountConditon() {
	BillService.Suffix.AddListener(binding.NewDataListener(func() {
		BillService.Load()
	}))
}

func dataListener() {
	date := time.Now().Format("200601")
	for {
		d := <-BillService.Condition.Date
		if len(d) == 4 {
			date = d + date[4:]
		} else if len(d) == 2 {
			date = date[:4] + d
		} else if len(d) == 1 {
			date = date[:4] + "0" + d
		} else {
			continue
		}

		start, err := time.Parse("20060102", date+"01")
		if err != nil {
			log.Panic(err)
		}

		end, err := time.Parse("20060102", date+"01")
		end = end.AddDate(0, 1, -1)
		if err != nil {
			log.Panic(err)
		}
		period := [2]time.Time{start, end}

		err = BillService.Condition.Period.Set(period)
		if err != nil {
			log.Panic(err)
		}

		BillService.UiEvent <- constant.Date
		log.Println("period update: ", period)
	}
}
