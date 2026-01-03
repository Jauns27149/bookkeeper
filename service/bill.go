package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/event"
	"bookkeeper/model"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

var _bill = &bill{}

type bill struct {
	_data     []model.Data
	aggregate model.Aggregate
	condition *model.Condition
	pref      fyne.Preferences
}

func (b *bill) run() {
	prefix, suffix := binding.NewString(), binding.NewString()
	prefix.Set(".*")
	suffix.Set(".*")
	n := time.Now()

	b.aggregate = model.Aggregate{
		Income:   binding.NewString(),
		Expenses: binding.NewString(),
		Budget:   binding.NewString(),
	}

	b.condition = &model.Condition{}
	b.condition.Account = map[string][]string{".*": {".*"}}
	b.condition.Perfix = binding.NewString()
	b.condition.Suffix = binding.NewString()
	b.condition.Start = time.Date(n.Year(), n.Month(), 1, 0, 0, 0, 0, time.UTC)
	b.condition.End = n

	b.pref = fyne.CurrentApp().Preferences()
	go func() {
		LoadBill()
		_bill.fetchConditionAccount()
		event.SetEventFunc(constant.LoadBill, LoadBill)
	}()

}

func LoadBill() {
	_bill.loadData()
	_bill.calculationAggregate()
	event.LaunchEvent(constant.BillRefresh)
}

func (b *bill) fetchConditionAccount() {
	accounts := b.pref.StringList(constant.Accounts)
	for _, a := range accounts {
		v := strings.Split(a, ":")
		if b.condition.Account[v[0]] == nil {
			b.condition.Account[v[0]] = append(b.condition.Account[v[0]], ".*")
		}
		b.condition.Account[v[0]] = append(b.condition.Account[v[0]], v[1])
	}

	event.LaunchEvent(constant.ConditionPrefixRefresh, constant.ConditionSuffixRefresh)
	log.Println("fetch condition account finished")
}

func GetCondition() *model.Condition {
	return _bill.condition
}

func GetAggregate() *model.Aggregate {
	return &_bill.aggregate
}

func Delete(id int) {
	_bill.delete(_bill._data[id])
}

func (b *bill) delete(deal model.Data) {
	key := deal.Date.Format(constant.OnlyMonth)
	list := b.pref.StringList(key)
	data := convert.RowsToDatas(list)
	for i, item := range data {
		if item == deal {
			b.pref.SetStringList(key, append(list[:i], list[i+1:]...))
			event.LaunchEvent(constant.LoadBill)
			log.Println("delete item successful, ", deal)
			break
		}
	}
}

func (b *bill) calculationAggregate() {
	accountMap := make(map[string]float64)
	for _, d := range b._data {
		accountMap[strings.Split(d.From.Name, ":")[0]] += d.From.Cost
		accountMap[strings.Split(d.To.Name, ":")[0]] += d.To.Cost
	}

	income := -accountMap[constant.Income]
	b.aggregate.Income.Set(strconv.FormatFloat(income, 'f', 2, 64))

	expenses := accountMap[constant.Expenses]
	b.aggregate.Expenses.Set(strconv.FormatFloat(expenses, 'f', 2, 64))

	budget := income*0.618 - expenses
	b.aggregate.Budget.Set(strconv.FormatFloat(budget, 'f', 2, 64))
}

func (b *bill) loadData() {
	t := b.condition.Start
	perfix, _ := b.condition.Perfix.Get()
	suffix, _ := b.condition.Suffix.Get()

	b._data = []model.Data{}
	for !t.After(b.condition.End) {
		d := b.pref.StringList(t.Format("2006-01"))
		_data := convert.RowsToDatas(d)
		_data = slices.DeleteFunc(_data, func(data model.Data) bool {
			return !regexp.MustCompile(perfix+suffix).MatchString(data.From.Name+data.To.Name) ||
				data.Date.Before(b.condition.Start) || data.Date.After(b.condition.End)

		})
		b._data = append(b._data, _data...)
		t = t.AddDate(0, 1, 0)
	}
	slices.SortFunc(b._data, func(x, y model.Data) int { return int(y.Date.Sub(x.Date)) })

	log.Println("load date finished, data size:", len(b._data))
}

func FetchData() []model.Data {
	return _bill._data
}

func DataByIndex(i int) model.Data {
	return _bill._data[i]
}
