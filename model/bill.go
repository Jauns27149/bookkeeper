package model

import (
	"time"

	"fyne.io/fyne/v2/data/binding"
)

type Aggregate struct {
	Income   binding.String
	Expenses binding.String
	Budget   binding.String
}

type Condition struct {
	Account map[string][]string
	Perfix  binding.String
	Suffix  binding.String
	Start   time.Time
	End     time.Time
}

type Data struct {
	Date     time.Time
	From     Account
	To       Account
	Terminal string
	Usage    string
}

type Account struct {
	Name string
	Cost float64
	Kind string
}
