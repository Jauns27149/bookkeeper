package util

import (
	"bookkeeper/model"
	"log"
	"regexp"
	"sort"
	"strconv"
)

func SortByDate(data []model.Deal) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Date.After(data[j].Date)
	})
}

func SortObjectsBySum(sum []string) {
	sort.Slice(sum, func(i, j int) bool {
		compile := regexp.MustCompile(`[0-9.]+$`)
		sumI, err := strconv.ParseFloat(compile.FindString(sum[i]), 64)
		if err != nil {
			log.Fatal(err)
		}
		sumJ, err := strconv.ParseFloat(compile.FindString(sum[j]), 64)
		if err != nil {
			log.Fatal(err)
		}

		return sumI > sumJ
	})
}
