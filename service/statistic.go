package service

import (
	"bookkeeper/constant"
	"bookkeeper/model"
	"math"
	"strings"
)

type Statistic struct {
	Deals      []model.Deal
	ExpenseMap map[string]float64
}

var StatisticService *Statistic

func newStatistic() *Statistic {
	data := &Statistic{
		Deals:      DataService.Deals,
		ExpenseMap: make(map[string]float64),
	}

	count := func(account string, cost float64) {
		if strings.Contains(account, constant.Expenses) {
			data.ExpenseMap[account] += cost
		}
	}
	for _, v := range DataService.Deals {
		count(v.AccountA, v.AccountAPay)
		count(v.AccountB, v.AccountBPay)
	}

	for k, v := range data.ExpenseMap {
		if v == 0 {
			delete(data.ExpenseMap, k)
			continue
		}
		data.ExpenseMap[k] = math.Round(v*100) / 100
	}

	return data
}

func StatisticRun() {
	StatisticService = newStatistic()
}
