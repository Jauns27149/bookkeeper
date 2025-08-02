package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	Period binding.Item[[2]time.Time]

	Pref      fyne.Preferences
	Deals     []model.Deal
	Income    binding.String
	Expense   binding.String
	Liability binding.String
	Payees    []string
	Accounts  binding.StringList

	AccountType binding.String

	Statements []model.Statement
	tally
}

type tally struct {
	Payee   binding.String
	Utilize binding.String
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
		compute(deal.Payment.Name, deal.Payment.Cost)
		compute(deal.Receiver.Name, deal.Receiver.Cost)
	}

	err := d.Expense.Set(strconv.FormatFloat(expense, 'f', 2, 64))
	if err != nil {
		log.Panicln(err.Error())
	}
	err = d.Income.Set(strconv.FormatFloat(-income, 'f', 2, 64))
	if err != nil {
		log.Panicln(err.Error())
	}
	err = d.Liability.Set(strconv.FormatFloat(-liability, 'f', 2, 64))
	if err != nil {
		log.Panicln(err.Error())
	}
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
		accountCount(v.Payment.Name)
		accountCount(v.Receiver.Name)
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
	DataService.count()
	DataService.payees()
	DataService.accounts()
	d.RefreshPage()
}

func (d *Data) RefreshPage() {
	fyne.CurrentApp().Driver().AllWindows()[0].Content().Refresh()
}

func (d *Data) Add(deal model.Deal) {
	key := deal.Date.Format("2006-01")
	datas := d.Pref.StringList(key)
	if len(datas) == 0 {
		periods := d.Pref.StringList(constant.Period)
		d.Pref.SetStringList(constant.Period, append(periods, key))
	}
	list := append(datas, convert.DealToString(deal))
	d.Pref.SetStringList(key, list)
	d.Refresh()
	TallyService.Finish <- constant.Bill
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
			log.Printf("remove deal %s success", target)
		}
	}
	d.Refresh()
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
	d.count()

	d.RefreshPage()
	log.Println("change current data success")
}

func (d *Data) LoadData() {
	period, err := d.Period.Get()
	if err != nil {
		log.Panicln(err)
	}
	start, end := period[0], period[1]

	pref, deals := d.Pref, make([]model.Deal, 0)
	for ; !start.After(end); start = start.AddDate(0, 1, 0) {
		for _, v := range pref.StringList(start.Format("2006-01")) {
			deal := convert.StringToDeal(v)
			account, err := d.AccountType.Get()
			if err != nil {
				log.Panicln(err.Error())
			}
			if !deal.Date.Before(start) && !deal.Date.After(end) && util.CheckAccount(deal, account) {
				deals = append(deals, deal)
			}
		}
	}
	util.SortByDate(deals)
	d.Deals = deals

	m := make(map[time.Time][]model.Deal)
	for _, v := range d.Deals {
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
	d.Statements = s

	d.count()
	d.RefreshPage()
	log.Println("load data finished")
}

func NewData() *Data {
	pref := fyne.CurrentApp().Preferences()
	now := time.Now()
	accountType := binding.NewString()
	err := accountType.Set(constant.All)
	if err != nil {
		log.Panicln(err.Error())
	}

	period := binding.NewItem(func(t [2]time.Time, t2 [2]time.Time) bool {
		return t == t2
	})
	err = period.Set([2]time.Time{
		time.Date(now.Year(), now.Month()-3, 1, 0, 0, 0, 0, time.Local),
		time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local),
	})
	if err != nil {
		log.Panicln(err)
	}

	data := &Data{
		Period:    period,
		Pref:      pref,
		Deals:     make([]model.Deal, 0),
		Income:    binding.NewString(),
		Expense:   binding.NewString(),
		Liability: binding.NewString(),
		Accounts:  binding.NewStringList(),

		AccountType: accountType,
	}
	data.LoadData()

	return data
}
