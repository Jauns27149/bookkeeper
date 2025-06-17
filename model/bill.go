package model

import "time"

type Deal struct {
	Date         time.Time
	Payee        string
	Usage        string
	AccountA     string
	AccountAPay  float64
	AccountAKind string
	AccountB     string
	AccountBPay  float64
	AccountBKind string
}
