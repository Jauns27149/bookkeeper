package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"bookkeeper/util"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"log"
	"regexp"
	"time"
)

type Bill struct {
	pref       fyne.Preferences
	Statements []model.Statement
	Head
	Condition

	DataEvent chan int
	UiEvent   chan int
}

type Condition struct {
	Period binding.Item[[2]time.Time]
	Prefix binding.String
	Suffix binding.String

	Date chan string
}

type Head struct {
	Income    binding.String
	Expense   binding.String
	Liability binding.String
	Budget    binding.String
}

func (b *Bill) Add(deal model.Deal) {
	key := deal.Date.Format(constant.YearMonth)
	list := b.pref.StringList(key)
	b.pref.SetStringList(key, append(list, convert.DealToRow(deal)))
	b.checkperiod(key)

	b.DataEvent <- constant.Load
	log.Println("add deal successful, ", deal)
}

func (b *Bill) Delete(deal model.Deal) {
	key := deal.Date.Format(constant.YearMonth)
	list := b.pref.StringList(key)
	for i, item := range list {
		if row := convert.DealToRow(deal); item == row {
			b.pref.SetStringList(key, append(list[:i], list[i+1:]...))
			b.DataEvent <- constant.Load
			log.Println("delete item successful, ", row)
			break
		}
	}
}

func (b *Bill) Load() {
	period := util.Period(b.Period)
	b.Statements = []model.Statement{}
	for _, month := range period {
		rows := b.pref.StringList(month)
		rows = b.assertCondition(rows)
		b.Statements = util.FillStatements(rows, b.Statements)
	}

	b.DataEvent <- constant.Count
	fmt.Println("load data successfully, ", len(b.Statements))
}

func (b *Bill) checkperiod(key string) {
	period := b.pref.StringList(constant.Period)
	for _, p := range period {
		if key == p {
			return
		}
	}
	b.pref.SetStringList(constant.Period, append(period, key))
	log.Printf("save bill key %s successful", key)
}

func (b *Bill) assertCondition(rows []string) (newRows []string) {
	account := util.AccountCombination(b.Prefix, b.Suffix)
	for _, v := range rows {
		if regexp.MustCompile(account).MatchString(v) {
			newRows = append(newRows, v)
		}
	}
	return
}

func NewBill() *Bill {
	return &Bill{
		pref: fyne.CurrentApp().Preferences(),
		Head: Head{
			Income:    binding.NewString(),
			Expense:   binding.NewString(),
			Liability: binding.NewString(),
			Budget:    binding.NewString(),
		},
		Condition: Condition{
			Period: util.CreatePeriod(),
			Prefix: binding.NewString(),
			Suffix: binding.NewString(),

			Date: make(chan string, 1),
		},
		Statements: []model.Statement{},
		DataEvent:  make(chan int, 1),
		UiEvent:    make(chan int, 1),
	}
}
