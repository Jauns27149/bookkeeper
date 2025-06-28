package util

import (
	"fmt"
	"testing"
)

func TestCountPay(t *testing.T) {
	m := make(map[string]float64)
	statement := "2025-05-26 马记永兰州牛肉面 午饭 Expenses:餐食 26.00 CNY Liabilities:信用卡 -26.00 CNY"
	CountPay(m, statement)
	fmt.Println(m)
}
