package util

import (
	"strconv"
	"time"
)

func Years(n int) (years []string) {
	start := time.Now().Year() - n
	years = make([]string, n*2+1)
	for i, _ := range years {
		years[i] = strconv.Itoa(start)
		start++
	}
	return
}

func Months(n int) (months []string) {
	for i := range 12 {
		months = append(months, strconv.Itoa(i+1))
	}
	for range n {
		months = append(months, months...)
	}
	return
}
