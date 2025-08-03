package util

import (
	"bookkeeper/model"
	"regexp"
)

func CountPay(accountMap map[string]float64, statement model.Deal) {
	for _, v := range []model.Account{statement.Payment, statement.Receiver} {
		accountMap[v.Name] += v.Cost
	}
}

func GroupAccountByPrefix(accounts []model.Bank) map[string][]model.Bank {
	reg := regexp.MustCompile(`^[^:]+`)
	accountMap := make(map[string][]model.Bank, len(accounts))
	for _, account := range accounts {
		prefix := reg.FindString(account.Account)
		if _, ok := accountMap[prefix]; !ok {
			accountMap[prefix] = []model.Bank{}
		}
		accountMap[prefix] = append(accountMap[prefix], account)
	}
	return accountMap
}

func CheckAccount(deal model.Deal, accountType string) bool {
	reg := regexp.MustCompile(accountType)
	return reg.MatchString("")
}

func LastAccount(s string) string {
	reg := regexp.MustCompile("[^:]+$")
	return reg.FindString(s)
}
