package convert

import (
	"bookkeeper/model"
	"cmp"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

func RowsToData(rows []string) (_data []model.Data) {
	for _, row := range rows {
		d := strings.Split(strings.TrimSpace(row), ",")
		if len(d) != 9 {
			log.Panicln(d)
		}
		date, _ := time.Parse(time.DateOnly, d[0])
		payA, _ := strconv.ParseFloat(d[4], 64)
		payB, _ := strconv.ParseFloat(d[7], 64)

		statement := model.Data{
			Date:     date,
			Terminal: d[1],
			Usage:    d[2],
			From:     model.Account{Name: d[3], Cost: payA, Kind: d[5]},
			To:       model.Account{Name: d[6], Cost: payB, Kind: d[8]},
		}
		if payA > 0 {
			statement.From, statement.To = statement.To, statement.From
		}

		_data = append(_data, statement)
	}

	return _data
}

func MapToAccounts(m map[string]map[string]float64) []model.AccountCategory {
	category := make([]model.AccountCategory, 0, len(m))
	for k, v := range m {
		account := &model.AccountCategory{Category: k}
		for k, c := range v {
			detail := model.AccountDetail{
				Name:   k,
				Amount: c,
			}

			account.AccountDetail = append(account.AccountDetail, detail)

			slices.SortFunc(account.AccountDetail, func(x, y model.AccountDetail) int {
				return int(x.Amount - y.Amount)
			})
		}

		category = append(category, *account)
		slices.SortFunc(category, func(x, y model.AccountCategory) int {
			return cmp.Compare(x.Category, y.Category)
		})
	}

	return category
}

func DataToRow(d model.Data) string {
	fmt.Println(d)
	rowData := []string{
		d.Date.Format(time.DateOnly),
		d.Terminal,
		d.Usage,
		d.From.Name,
		strconv.FormatFloat(d.From.Cost, 'f', 2, 64), d.From.Kind,
		d.To.Name,
		strconv.FormatFloat(d.To.Cost, 'f', 2, 64), d.To.Kind,
	}

	row := strings.Join(rowData, ",")
	log.Println("data to row finished, row:", row)
	return row
}

func StringToFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Panicln("convert failed," + err.Error())
	}

	return v
}
