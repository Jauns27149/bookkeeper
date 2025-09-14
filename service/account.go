package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"strings"
)

type Account struct {
	pref fyne.Preferences
}

func (r *Account) AccountCollection() []model.Bank {
	periods := r.pref.StringList(constant.Period)
	accountMap := make(map[string]float64)

	for _, period := range periods {
		statements := r.pref.StringList(period)
		for _, statement := range statements {
			deal := convert.StringToDeal(statement)
			util.CountPay(accountMap, deal)
		}
	}

	return convert.MapToBank(accountMap)
}

func (r *Account) AccountMap() (accountMap map[string][]string) {
	accountMap = map[string][]string{constant.All: {constant.All}}
	accounts := r.pref.StringList(constant.Accounts)
	for _, v := range accounts {
		value := strings.Split(v, ":")
		if _,ok :=accountMap[value[0]];!ok{
			accountMap[value[0]]=[]string{constant.All}
		}
		accountMap[value[0]] = append(accountMap[value[0]], value[1])
	}

	return
}

func NewAccount() *Account {
	return &Account{
		pref: fyne.CurrentApp().Preferences(),
	}
}
