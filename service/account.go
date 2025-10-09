package service

import (
	"bookkeeper/app"
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"log"
	"strings"

	"fyne.io/fyne/v2"
)

var accounts = &Accounts{
	pref: app.Preferences(),
}

var accountsFlag = make(chan struct{})

type Accounts struct {
	pref     fyne.Preferences
	Accounts []model.AccountCategory
}

func GetAccounts() *Accounts {
	<-accountsFlag
	return accounts
}

func init() {
	go func() {
		dataAll := make(map[string]map[string]float64, 5)
		period := accounts.pref.StringList(constant.Period)
		for _, p := range period {
			rows := accounts.pref.StringList(p)
			data := convert.RowsToDatas(rows)
			for _, d := range data {
				for _, v := range []model.Account{d.From, d.To} {
					names := strings.Split(v.Name, ":")
					if _, ok := dataAll[names[0]]; !ok {
						dataAll[names[0]] = make(map[string]float64)
					}
					dataAll[names[0]][names[1]] += v.Cost
				}
			}
		}

		for key, value := range dataAll {
			total := 0.0
			for _, v := range value {
				total += v
			}

			dataAll[key][""] = total
		}

		accounts.Accounts = convert.MapToAccounts(dataAll)
		close(accountsFlag)
		log.Println("load accounts data finished, size: ", len(accounts.Accounts))

	}()
}
