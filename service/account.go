package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/event"
	"bookkeeper/model"
	"log"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
)

var accounts = &Accounts{}

var accountsFlag = make(chan struct{})

type Accounts struct {
	pref     fyne.Preferences
	Accounts []model.AccountCategory
}

func GetAccounts() *Accounts {
	<-accountsFlag
	return accounts
}

func (a *Accounts) save(account string) {
	accounts := a.pref.StringList(constant.Accounts)
	accounts = append(accounts, account)
	a.pref.SetStringList(constant.Accounts, accounts)

	a.chnageAccount()
}

func (a *Accounts) rename(oldName, newName string) {
	period := a.pref.StringList(constant.Period)
	for _, p := range period {
		rows := a.pref.StringList(p)
		for _, row := range rows {
			if strings.Contains(row, oldName) {
				strings.ReplaceAll(row, oldName, newName)
			}
		}
		a.pref.SetStringList(p, rows)
	}
}

func (a *Accounts) run() {
	a.pref = fyne.CurrentApp().Preferences()
	go func() {
		dataAll := make(map[string]map[string]float64, 5)
		for _, v := range a.pref.StringList(constant.Accounts) {
			strs := strings.Split(v, ":")
			if _, ok := dataAll[strs[0]]; !ok {
				dataAll[strs[0]] = make(map[string]float64)
			}
			dataAll[strs[0]][strs[1]] = 0

		}

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

		go a.chnageAccount()
		log.Println("load accounts data finished, size: ", len(accounts.Accounts))

	}()
}

func (a *Accounts) chnageAccount() {
	m := _bill.condition.Account
	for _, v := range a.Accounts {
		for _, vv := range v.AccountDetail {
			name := vv.Name
			if name == "" {
				name = ".*"
			}
			if !slices.Contains(m[v.Category], name) {
				m[v.Category] = append(m[v.Category], name)
			}
		}
	}
	_bill.condition.Account = m

	event.LaunchEvent(constant.ConditionPrefixRefresh)
}
