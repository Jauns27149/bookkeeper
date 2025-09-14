package convert

import (
	"bookkeeper/constant"
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

func DealToRow(deal model.Deal) string {
	rowData := []string{
		deal.Date.Format(time.DateOnly),
		deal.Payee,
		deal.Usage,
		deal.Payment.Name,
		strconv.FormatFloat(deal.Payment.Cost, 'f', 2, 64), deal.Payment.Kind,
		deal.Receiver.Name,
		strconv.FormatFloat(deal.Receiver.Cost, 'f', 2, 64), deal.Receiver.Kind,
	}
	return strings.Join(rowData, constant.Comma)
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
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
	d := strings.Split(strings.TrimSpace(row), constant.Comma)
	if len(d) != 9 {
		log.Panicln(d)
	}
	date, _ := time.Parse(time.DateOnly, d[0])
	payA, _ := strconv.ParseFloat(d[4], 64)
	payB, _ := strconv.ParseFloat(d[7], 64)
	return model.Deal{
		Date:     date,
		Payee:    d[1],
		Usage:    d[2],
		Payment:  model.Account{Name: d[3], Cost: payA, Kind: d[5]},
		Receiver: model.Account{Name: d[6], Cost: payB, Kind: d[8]},
	}
}
