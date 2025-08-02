package model

import "time"

type Deal struct {
	Date     time.Time
	Payee    string
	Usage    string
	Payment  Account
	Receiver Account
}

type Account struct {
	Name string
	Cost float64
	Kind string
}

type Bank struct {
	Account string
	Amount  float64
}

type Statement struct {
	Date  time.Time
	Deals []Deal
}
