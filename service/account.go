package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
)

type Account struct {
	pref fyne.Preferences
}

func (a *Account) AccountCollection() []model.Bank {
	periods := a.pref.StringList(constant.Period)
	accountMap := make(map[string]float64)

	for _, period := range periods {
		statements := a.pref.StringList(period)
		for _, statement := range statements {
			util.CountPay(accountMap, statement)
		}
	}

	return convert.MapToBank(accountMap)
}

func NewAccount() *Account {
	return &Account{
		pref: fyne.CurrentApp().Preferences(),
	}
}
