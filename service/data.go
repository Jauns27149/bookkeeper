package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

var DataService *Data

type Data struct {
	Pref      fyne.Preferences
	Deals     []model.Deal
	Income    binding.Float
	Expense   binding.Float
	Liability binding.Float
	Payees    []string
	Accounts  binding.StringList
}

func (d *Data) count() {
	var income, expense, liability float64
	compute := func(account string, payment float64) {
		key := strings.ToLower(strings.Split(account, ":")[0])
		switch key {
		case "income":
			income += payment
		case "expenses":
			expense += payment
		case "liabilities":
			liability += payment
		}
	}

	for _, deal := range d.Deals {
		compute(deal.AccountA, deal.AccountAPay)
		compute(deal.AccountB, deal.AccountBPay)
	}

	d.Expense.Set(math.Round(expense*100) / 100)
	d.Income.Set(-math.Round(income*100) / 100)
	d.Liability.Set(-math.Round(liability*100) / 100)
}

func (d *Data) payees() {
	m := make(map[string]int)
	for _, v := range d.Deals {
		m[v.Payee]++
	}
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	sort.Slice(s, func(i, j int) bool {
		return m[s[i]] > m[s[j]]
	})
	d.Payees = s[:min(8, len(s))]
}

func (d *Data) accounts() {
	accounts := d.Pref.StringList(constant.Accounts)

	m := make(map[string]int)
	for _, v := range accounts {
		m[v] = 0
	}
	accountCount := func(account string) {
		m[account]++
	}
	for _, v := range d.Deals {
		accountCount(v.AccountA)
		accountCount(v.AccountB)
	}

	sort.Slice(accounts, func(i, j int) bool {
		return m[accounts[i]] > m[accounts[j]]
	})
	err := d.Accounts.Set(accounts)
	if err != nil {
		log.Println(err.Error())
	}
}

func (d *Data) Refresh() {
	d.GetDeals("")
	DataService.count()
	DataService.payees()
	DataService.accounts()
	d.RefreshPage()
}

func (d *Data) RefreshPage() {
	fyne.CurrentApp().Driver().AllWindows()[0].Content().Refresh()
}

func (d *Data) Save(deal model.Deal) {
	key := deal.Date.Format("2006-01")
	list := append(d.Pref.StringList(key), convert.DealToString(deal))
	d.Pref.SetStringList(key, list)
	d.Refresh()
	log.Printf("sava deal success, month %v, amount %v", key, len(list))
}

func (d *Data) GetDeals(key string) {
	//span := time.Now().AddDate(0, -1, 0).Format("2006-01")
	span := time.Now().Format("2006-01")
	data := d.Pref.StringList(span)
	currenData := make([]model.Deal, 0, len(data))

	for _, v := range data {
		if key != "" && !strings.Contains(v, key) {
			continue
		}
		bill := convert.StringToDeal(v)
		currenData = append(currenData, bill)
	}
	util.SortByDate(currenData)
	d.Deals = currenData
}

func (d *Data) RemoveDeal(deal model.Deal) {
	target := convert.DealToString(deal)
	key := deal.Date.Format("2006-01")
	items := d.Pref.StringList(key)
	for i, vv := range items {
		if vv == target {
			d.Pref.SetStringList(key, append(items[:i], items[i+1:]...))
		}
	}
	d.Refresh()
	log.Printf("remove deal %s success", target)
}

func (d *Data) ChangeDataByPeriod(start time.Time, end time.Time) {
	temp, _ := time.Parse(time.DateTime, start.Format(time.DateTime))
	newData := make([]model.Deal, 0)
	start = start.AddDate(0, 0, -1)
	end = end.AddDate(0, 0, 1)
	for !temp.After(end) {
		datas := d.Pref.StringList(temp.Format("2006-01"))
		for _, data := range datas {
			deal := convert.StringToDeal(data)
			if deal.Date.After(start) && deal.Date.Before(end) {
				newData = append(newData, deal)
			}
		}
		temp = temp.AddDate(0, 1, 0)
	}

	util.SortByDate(newData)
	d.Deals = newData
	log.Println("change current data success")
}

func (d *Data) FilterDataByAccount(account string) {
	newData := make([]model.Deal, 0)
	for _, deal := range d.Deals {
		if strings.Contains(deal.AccountA, account) || strings.Contains(deal.AccountB, account) {
			newData = append(newData, deal)
		}
	}
	util.SortByDate(newData)
	d.Deals = newData
	d.RefreshPage()
	log.Println("filter data success")
}

func DataRun() {
	pref := fyne.CurrentApp().Preferences()
	DataService = &Data{
		Pref:      pref,
		Deals:     make([]model.Deal, 0),
		Income:    binding.NewFloat(),
		Expense:   binding.NewFloat(),
		Liability: binding.NewFloat(),
		Accounts:  binding.NewStringList(),
	}

	DataService.GetDeals("")
	DataService.Refresh()
}
