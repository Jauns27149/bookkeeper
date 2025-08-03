package util

import (
	"fyne.io/fyne/v2/data/binding"
	"log"
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

func CreatePeriod() binding.Item[[2]time.Time] {
	period := binding.NewItem(func(t, tt [2]time.Time) bool {
		return t[1].Equal(tt[1])
	})
	now, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
	if err != nil {
		log.Panic(err)
	}
	err = period.Set([2]time.Time{
		now.AddDate(0, -3, -now.Day()+1),
		now.AddDate(0, 1, -now.Day())})
	if err != nil {
		log.Panic(err)
	}
	return period
}

func Period(period binding.Item[[2]time.Time]) []string {
	start, end := timeForPeriod(period)
	end.AddDate(0, 1, 0)
	result := make([]string, 0)
	for !start.After(end) {
		result = append(result, start.Format("2006-01"))
		start = start.AddDate(0, 1, 0)
	}
	return result
}

func timeForPeriod(period binding.Item[[2]time.Time]) (time.Time, time.Time) {
	p, err := period.Get()
	if err != nil {
		log.Panicln(err)
	}
	return p[0], p[1]
}
