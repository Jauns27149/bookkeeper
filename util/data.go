package util

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"fyne.io/fyne/v2/data/binding"
	"log"
	"strconv"
	"strings"
	"time"
)

func FillStatements(rows []string, statements []model.Statement) []model.Statement {
	for _, row := range rows {
		deal := convert.StringToDeal(row)
		statements = addStatements(statements, deal)
	}
	return statements
}

func addStatements(statements []model.Statement, deal model.Deal) []model.Statement {
	for i, statement := range statements {
		if statement.Date.Equal(deal.Date) {
			statements[i].Deals = append(statement.Deals, deal)
			return statements
		}
		if deal.Date.After(statements[i].Date) {
			temp := model.Statement{Date: deal.Date, Deals: []model.Deal{deal}}
			return append(statements[:i], append([]model.Statement{temp}, statements[i:]...)...)
		}
	}
	return append(statements, model.Statement{Date: deal.Date, Deals: []model.Deal{deal}})
}

func CountHead(m map[string]float64, deal model.Deal) {
	for _, v := range []model.Account{deal.Payment, deal.Receiver} {
		key := strings.Split(v.Name, ":")[0]
		m[key] += v.Cost
	}
}

func SetValue(bind binding.String, value float64) {
	err := bind.Set(strconv.FormatFloat(value, 'f', 2, 64))
	if err != nil {
		log.Panicln(err)
	}
}

func CountData(statements []model.Statement) map[string]float64 {
	m := make(map[string]float64)
	for _, statement := range statements {
		for _, v := range statement.Deals {
			CountHead(m, v)
		}
	}
	return m
}

func BindSting(bind binding.String) string {
	value, err := bind.Get()
	if err != nil {
		log.Panicln(err)
	}
	return value
}

func Account(bind model.BindAccount, b binding.String, positive bool) model.Account {
	name := BindSting(bind.Prefix) + ":" + BindSting(bind.Account)
	cost, err := strconv.ParseFloat(BindSting(b), 64)
	if !positive {
		cost = -cost
	}
	if err != nil {
		log.Panicln(err)
	}
	return model.Account{
		Name: name,
		Kind: "CNY",
		Cost: cost,
	}
}

func BindTime(date string, bind binding.String) time.Time {
	t, err := time.Parse(constant.FyneDate, BindSting(bind))
	if err != nil {
		log.Panicln(err)
	}
	return t
}
