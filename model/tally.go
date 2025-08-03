package model

import "fyne.io/fyne/v2/data/binding"

type SubAccount struct {
	Prefix   string
	Suffixes []string
}

type Record struct {
	Cost     binding.String
	Receiver binding.String
	Usage    binding.String
	From     BindAccount
	To       BindAccount
	Date     binding.String
}

type BindAccount struct {
	Prefix  binding.String
	Account binding.String
}
