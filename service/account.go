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

	a.changeAccount()
}

func (a *Accounts) Rename(oldName, newName string) {
	list := a.pref.StringList(constant.Accounts)
	for i, account := range list {
		if strings.Contains(account, oldName) {
			list[i] = strings.ReplaceAll(account, oldName, newName)
		}
	}
	a.pref.SetStringList(constant.Accounts, list)

	period := a.pref.StringList(constant.Period)
	for _, p := range period {
		rows := a.pref.StringList(p)
		for i, row := range rows {
			if strings.Contains(row, oldName) {
				rows[i] = strings.ReplaceAll(row, oldName, newName)
			}
		}
		a.pref.SetStringList(p, rows)
	}

	log.Printf("rename finished,%v-->%v\n", oldName, newName)
}

func (a *Accounts) run() {
	a.pref = fyne.CurrentApp().Preferences()
	go func() {
		dataAll := make(map[string]map[string]float64, 5)
		for _, v := range a.pref.StringList(constant.Accounts) {
			row := strings.Split(v, ":")
			if _, ok := dataAll[row[0]]; !ok {
				dataAll[row[0]] = make(map[string]float64)
			}
			dataAll[row[0]][row[1]] = 0
		}

		period := accounts.pref.StringList(constant.Period)
		for _, p := range period {
			rows := accounts.pref.StringList(p)
			data := convert.RowsToData(rows)
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

		go a.changeAccount()
		log.Println("load accounts data finished, size: ", len(accounts.Accounts))

	}()
}

func (a *Accounts) changeAccount() {
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
