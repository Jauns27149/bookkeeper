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
		b.Statements = util.FillStatements(rows, b.Statements)
	}

	b.DataEvent <- constant.Count
	fmt.Println("load data successfully, ", len(b.Statements))
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
		},
		Statements: []model.Statement{},
		DataEvent:  make(chan int, 1),
		UiEvent:    make(chan int, 1),
	}
}
