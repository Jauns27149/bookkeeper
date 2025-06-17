package main

import (
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	log.Println("账本开始启动...")
	a := app.NewWithID("bookkeeper")
	w := a.NewWindow("bookkeeper")

	//initData()

	w.Resize(fyne.NewSize(800, 600))
	bill := widget.NewButton("账单", nil)
	statistic := widget.NewButton("统计", nil)
	statisticContent := widget.NewLabel("暂时没有没有更新")
	billService := service.NewBill()

	columns := container.NewGridWithColumns(2)
	columns.Add(bill)
	columns.Add(statistic)
	content := container.NewBorder(nil, columns, nil, nil, billService.Content, statisticContent)
	statisticContent.Hide()

	bill.Disable()
	bill.OnTapped = func() {
		c := w.Content().(*fyne.Container)
		bill.Disable()
		statistic.Enable()
		statisticContent.Hide()
		billService.Content.Show()
		c.Refresh()
	}

	statistic.OnTapped = func() {
		c := w.Content().(*fyne.Container)
		statistic.Disable()
		bill.Enable()
		billService.Content.Hide()
		statisticContent.Show()
		c.Refresh()
	}

	w.SetContent(content)
	w.ShowAndRun()
}

func initData() {
	// assets.bean  equity.bean  expenses.bean  income.bean  liabilities.bean
	preferences := fyne.CurrentApp().Preferences()

	files, _ := os.ReadDir("data")
	dataMap := make(map[string][]string)
	accountMap := make(map[string]bool)

	for _, f := range files {
		name := f.Name()
		dataByte, _ := os.ReadFile("data/" + name)
		data := strings.TrimSpace(string(dataByte))
		data = strings.ReplaceAll(data, " * ", " ")
		data = strings.ReplaceAll(data, "\r\n", "\n")
		reg := regexp.MustCompile(`\n{2,}`)
		if reg.MatchString(data) {
			data = reg.ReplaceAllString(data, "===")
		}

		reg = regexp.MustCompile(`:.*:`)
		if reg.MatchString(data) {
			data = reg.ReplaceAllString(data, ":")
		}

		reg = regexp.MustCompile(`\n`)
		if reg.MatchString(data) {
			data = reg.ReplaceAllString(data, " ")
		}

		reg = regexp.MustCompile(`[A-Za-z]*:\p{Han}*`)
		for _, v := range reg.FindAllString(data, -1) {
			if v = strings.TrimSpace(v); !accountMap[v] {
				accountMap[v] = true
			}
		}

		reg = regexp.MustCompile(`("")`)
		data = reg.ReplaceAllString(data, "---")

		data = strings.ReplaceAll(data, `"`, "")

		reg = regexp.MustCompile(`\s+`)
		data = reg.ReplaceAllString(data, " ")

		deal := strings.Split(data, "===")
		dataMap[name[:len(name)-5]] = deal
	}
	for k, v := range dataMap {
		preferences.RemoveValue(k)
		preferences.SetStringList(k, v)
	}
	preferences.RemoveValue("2025-06")

	accounts := make([]string, 0, len(accountMap))
	for k, _ := range accountMap {
		accounts = append(accounts, k)
	}
	preferences.SetStringList("accounts", accounts)

}
