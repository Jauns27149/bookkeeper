package model

type AccountCategory struct {
	Category      string
	AccountDetail []AccountDetail
}

type AccountDetail struct {
	Name   string
	Amount float64
}
