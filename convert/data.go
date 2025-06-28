package convert

import (
	"bookkeeper/model"
	"github.com/wcharczuk/go-chart"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func MapToChartValues(Map map[string]float64) []chart.Value {
	values := make([]chart.Value, 0, len(Map))
	for k, v := range Map {
		value := chart.Value{Value: v, Label: k}
		values = append(values, value)
	}
	return values
}

func DealToString(deal model.Deal) string {
	rowData := []string{
		deal.Date.Format(time.DateOnly),
		deal.Payee,
		deal.Usage,
		deal.AccountA,
		strconv.FormatFloat(deal.AccountAPay, 'f', -1, 64), deal.AccountAKind,
		deal.AccountB,
		strconv.FormatFloat(deal.AccountBPay, 'f', -1, 64), deal.AccountBKind,
	}
	return strings.Join(rowData, " ")
}

func MapToBank(accountMap map[string]float64) []model.Bank {
	bank := make([]model.Bank, 0, len(accountMap))
	for k, v := range accountMap {
		bank = append(bank, model.Bank{Account: k, Amount: v})
	}
	sort.Slice(bank, func(i, j int) bool {
		return bank[i].Amount > bank[j].Amount
	})
	return bank
}

func StringToDeal(row string) model.Deal {
	d := strings.Split(row, " ")
	if len(d) < 9 {
		log.Panicln(d)
	}
	date, _ := time.Parse(time.DateOnly, d[0])
	payA, _ := strconv.ParseFloat(d[4], 64)
	payB, _ := strconv.ParseFloat(d[7], 64)
	return model.Deal{
		Date:         date,
		Payee:        d[1],
		Usage:        d[2],
		AccountA:     d[3],
		AccountAPay:  payA,
		AccountAKind: d[5],
		AccountB:     d[6],
		AccountBPay:  payB,
		AccountBKind: d[8],
	}
}
