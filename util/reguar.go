package util

import (
	"bookkeeper/model"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func CountPay(accountMap map[string]float64, statement string) {
	reg := regexp.MustCompile(`[^ ]+:[^ ]+ [0-9.-]+`)
	all := reg.FindAllString(statement, -1)
	for _, v := range all {
		banks := strings.Split(v, " ")
		account := banks[0]
		cost, err := strconv.ParseFloat(banks[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		accountMap[account] += cost
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
