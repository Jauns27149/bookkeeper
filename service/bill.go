package service

import (
	"bookkeeper/constant"
	"bookkeeper/model"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Bill struct {
	Content       fyne.CanvasObject
	head          *fyne.Container
	record        fyne.CanvasObject
	data          []model.Deal
	pref          fyne.Preferences
	accounts      map[string]int
	startEntry    *widget.DateEntry
	endEntry      *widget.DateEntry
	scroll        *container.Scroll
	accountList   []string
	accountSelect []*widget.Select
	headButton    []*widget.Button
}

func NewBill() (b *Bill) {
	b = new(Bill)
	b.pref = fyne.CurrentApp().Preferences()
	b.accounts = make(map[string]int)
	for _, d := range b.pref.StringList("accounts") {
		b.accounts[d] = 0
	}
	b.accountList = b.pref.StringList(constant.Accounts)

	b.loadData()
	b.createHead()
	b.createContent()
	return
}

func (b *Bill) createContent() {
	buttons := make([]*widget.Button, len(b.data))
	updates := make([]*widget.Button, len(b.data))
	//buttons := make([]*widget.Button, 0, len(b.data))
	list := widget.NewList(
		func() int {
			//buttons = make([]*widget.Button, len(b.data))
			if len(buttons) < len(b.data) {
				buttons = append(buttons, nil)
				updates = append(updates, nil)
			}
			return len(b.data)
		},

		func() fyne.CanvasObject {
			//return widget.NewLabel("")
			items := make([]fyne.CanvasObject, 3)
			for i := range 3 {
				items[i] = canvas.NewText("", nil)
			}
			h := container.NewVBox(items...)

			items = make([]fyne.CanvasObject, 5)
			items[0] = h
			items[1] = layout.NewSpacer()
			items[2] = canvas.NewText("", nil)
			items[3] = widget.NewButton("删除", nil)
			items[4] = widget.NewButton("更新", nil)
			v := container.NewHBox(items...)
			return v
		},

		func(id widget.ListItemID, object fyne.CanvasObject) {
			v := object.(*fyne.Container)
			left := v.Objects[0].(*fyne.Container)
			pay := v.Objects[2].(*canvas.Text)
			button := v.Objects[3].(*widget.Button)
			update := v.Objects[4].(*widget.Button)

			deal := b.data[id]
			left.Objects[0].(*canvas.Text).Text = deal.Usage
			left.Objects[1].(*canvas.Text).Text = deal.Payee
			head, tail := "", ""
			if deal.AccountAPay < 0 {
				head, tail = deal.AccountA, deal.AccountB
			} else {
				head, tail = deal.AccountB, deal.AccountA
			}
			text := fmt.Sprintf("%v %v->%v", deal.Date.Format(time.DateOnly), head, tail)
			left.Objects[2].(*canvas.Text).Text = text
			pay.Text = strconv.FormatFloat(deal.AccountAPay, 'f', 2, 64)
			button.Importance = widget.WarningImportance
			button.Hide()
			update.Hide()
			button.OnTapped = func() {
				b.data = append(b.data[:id], b.data[id+1:]...)
				b.Content.Refresh()
				value := make([]string, len(b.data))
				for i, d := range b.data {
					// "2025-05-17 美团 麦蹈中式健康菜外卖 Expenses:餐食 19.90 CNY Liabilities:信用卡 -19.90 CNY"
					strs := []string{
						d.Date.Format(time.DateOnly), d.Payee, d.Usage,
						d.AccountA, strconv.FormatFloat(d.AccountAPay, 'f', -1, 64), d.AccountAKind,
						d.AccountB, strconv.FormatFloat(d.AccountBPay, 'f', -1, 64), d.AccountBKind,
					}
					str := strings.Join(strs, " ")
					value[i] = str
				}
				b.pref.SetStringList(time.Now().Format("2006-01"), value)
				//button.Hide()
			}

			update.OnTapped = func() {
				//update.Hide()
				//buttons[id].TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
				data := b.data[id]
				buttons[id].OnTapped()

				box := b.record.(*fyne.Container)
				items := box.Objects
				dateEntry := items[0].(*widget.DateEntry)
				dateEntry.SetDate(&data.Date)
				selectEntry := items[1].(*widget.SelectEntry)
				selectEntry.SetText(data.Payee)
				entry := items[2].(*widget.Entry)
				entry.SetText(data.Usage)
				objects := items[3].(*fyne.Container).Objects
				w := objects[0].(*widget.Select)
				w2 := objects[1].(*widget.Entry)
				w.SetSelected(data.AccountA)
				w2.SetText(strconv.FormatFloat(data.AccountAPay, 'f', 2, 64))

				objects = items[4].(*fyne.Container).Objects
				w = objects[0].(*widget.Select)
				w2 = objects[1].(*widget.Entry)
				w.SetSelected(data.AccountB)
				w2.SetText(strconv.FormatFloat(data.AccountBPay, 'f', 2, 64))
				b.record.Show()
			}
			//buttons = append(buttons, button)
			buttons[id] = button
			updates[id] = update
			log.Println("add delete button successful")
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		log.Println("selected delete button")
		buttons[id].Show()
		updates[id].Show()
		//b.Content.Refresh()
	}
	content := container.NewBorder(b.head, nil, nil, nil, container.NewBorder(b.record, nil, nil, nil, list))
	b.Content = content
}

func (b *Bill) createHead() {
	head := container.NewGridWithColumns(1)
	b.head = head

	headButton := container.NewGridWithColumns(3)
	countMap := map[string]float64{constant.Expenses: 0, constant.Income: 0, constant.Liabilities: 0}
	for _, v := range b.data {
		temp := map[string]float64{v.AccountA: v.AccountAPay, v.AccountB: v.AccountBPay}
		for k, vv := range temp {
			switch kk := strings.Split(k, ":")[0]; kk {
			case constant.Expenses, constant.Liabilities, constant.Income:
				countMap[kk] = countMap[kk] + vv
			}
		}
	}
	b.headButton = make([]*widget.Button, 3)
	for i, v := range []string{"收入", "支出", "负债"} {
		text := v
		var count float64
		switch i {
		case 0:
			count = -countMap[constant.Income]
		case 1:
			count = countMap[constant.Expenses]
		case 2:
			count = -countMap[constant.Liabilities]
		}
		text = fmt.Sprintf("%v:%s", text, strconv.FormatFloat(count, 'f', 2, 64))
		button := widget.NewButton(text, nil)
		headButton.Add(button)
		b.headButton[i] = button
	}
	head.Add(headButton)
	//entry := widget.NewSelect([]string{"Expenses"}, nil)
	//entry.SetSelected("Expenses")
	//account := container.NewHBox(widget.NewLabel("账户类型："), entry)
	//head.Add(account)

	createButton := widget.NewButton("创建账号", nil)
	head.Add(createButton)

	accountSelect := widget.NewSelect(constant.AccountsList, nil)
	accountSelect.SetSelected(constant.Expenses)
	accountEntry := widget.NewEntry()
	accountGrid := container.NewGridWithColumns(2, accountSelect, accountEntry)
	accountGrid.Hide()
	head.Add(accountGrid)

	accountEntry.OnSubmitted = func(s string) {
		newAccount := fmt.Sprintf("%s:%s", accountSelect.Selected, accountEntry.Text)
		accountGrid.Hide()
		b.accountList = append(b.accountList, newAccount)
		b.pref.SetStringList(constant.Accounts, b.accountList)
		for i, v := range b.accountSelect {
			v.SetOptions(b.accountList)
			v.SetSelectedIndex(i)
		}
	}

	createButton.OnTapped = func() {
		if accountGrid.Hidden {
			accountGrid.Show()
			return
		}
		accountGrid.Hide()
	}

	b.createRecord()

	//now := time.Now()
	//var start time.Time
	//for i, v := range []string{"开始时间：", "结束时间："} {
	//	dateEntry := widget.NewDateEntry()
	//	border := container.NewBorder(nil, nil, widget.NewLabel(v), nil, dateEntry)
	//	head.Add(border)
	//	switch i {
	//	case 0:
	//		s := now.Format("2006-01-")
	//		start, _ = time.Parse("2006-01-02", s+"01")
	//		dateEntry.SetDate(&start)
	//		b.startEntry = dateEntry
	//	case 1:
	//		dateEntry.SetDate(&now)
	//		dateEntry.OnSubmitted = func(s string) {
	//			b.createContent()
	//		}
	//		b.endEntry = dateEntry
	//	}

	//}
}

func (b *Bill) createRecord() {
	record := widget.NewButton("+记账", nil)
	b.head.Add(record)
	box := container.NewVBox()

	dateEntry := widget.NewDateEntry()
	now := time.Now()
	dateEntry.SetDate(&now)
	box.Add(dateEntry)

	m := make(map[string]int)
	for _, v := range b.data {
		m[v.Payee]++
	}
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	sort.Slice(s, func(i, j int) bool {
		return m[s[i]] > m[s[j]]
	})
	s = s[:min(5, len(s))]
	payee := widget.NewSelectEntry(s)
	payee.SetPlaceHolder("收款人/商户/收入来源渠道")
	if len(s) > 0 {
		payee.SetText(s[0])
	}
	box.Add(payee)

	detail := widget.NewEntry()
	detail.SetPlaceHolder("详细描述，记录细节")
	box.Add(detail)

	preferences := fyne.CurrentApp().Preferences()
	clear(m)
	for _, v := range b.data {
		m[v.AccountA]++
		m[v.AccountB]++
	}

	b.accountSelect = make([]*widget.Select, 2)
	amounts := make([]*widget.Entry, 2)
	for i := range 2 {
		b.accountSelect[i] = widget.NewSelect(b.accountList, nil)
		account := b.accountSelect[i]
		if b.accountList != nil && len(b.accountList) > i {
			account.SetSelectedIndex(i)
		}
		amounts[i] = widget.NewEntry()
		amount := amounts[i]
		amount.SetPlaceHolder("金额")
		amount.OnSubmitted = func(s string) {
			if v := amounts[0].Text; v == "" {
				value := amounts[1].Text
				amounts[0].SetText("-" + value)
			}
			if v := amounts[1].Text; v == "" {
				value := amounts[0].Text
				amounts[1].SetText("-" + value)
			}

			deal := make([]string, 9)
			deal[0] = dateEntry.Date.Format("2006-01-02")
			deal[1] = payee.Text
			deal[2] = detail.Text
			deal[3] = b.accountSelect[0].Selected
			deal[4] = amounts[0].Text
			deal[5] = "CNY"
			deal[6] = b.accountSelect[1].Selected
			deal[7] = amounts[1].Text
			deal[8] = "CNY"

			amounts[0].SetText("")
			amounts[1].SetText("")

			list := preferences.StringList(now.Format("2006-01"))
			list = append(list, strings.Join(deal, " "))
			preferences.SetStringList(now.Format("2006-01"), list)
			box.Hide()
			b.loadData()
			b.Content.Refresh()
		}
		acc := container.NewGridWithColumns(2)
		acc.Add(account)
		acc.Add(amount)
		box.Add(acc)
	}
	// 模板 "2024-11-23 交通 地铁 Assets:储蓄卡 -10.00 CNY Expenses:公共交通 10.00 CNY"
	box.Hide()
	record.OnTapped = func() {
		if box.Hidden {
			box.Show()
		} else {
			box.Hide()
			b.Content.Refresh()
			fyne.CurrentApp().Driver().AllWindows()[0].Content().Refresh()
		}
	}
	b.record = box
}

func (b *Bill) loadData() {
	end := time.Now()
	start, err := time.Parse(time.DateOnly, end.Format("2006-01-")+"01")
	if err != nil {
		log.Println(err)
	}

	data := make([]string, 0)
	for start.Month() <= end.Month() {
		key := start.Format("2006-01")
		list := b.pref.StringList(key)
		data = append(data, list...)
		start = start.AddDate(0, 1, 0)
	}

	bills := make([]model.Deal, 0, len(data))
	for _, v := range data {
		d := strings.Split(v, " ")
		if len(d) < 9 {
			log.Panicln(d)
		}
		date, _ := time.Parse(time.DateOnly, d[0])
		payA, _ := strconv.ParseFloat(d[4], 64)
		payB, _ := strconv.ParseFloat(d[7], 64)
		bill := model.Deal{
			Date:         date,
			Payee:        d[1],
			Usage:        d[2],
			AccountA:     d[3],
			AccountAPay:  payA,
			AccountAKind: d[5],
			AccountB:     d[6],
			AccountBPay:  payB,
			AccountBKind: d[8]}
		bills = append(bills, bill)
		b.accounts[bill.AccountA]++
		b.accounts[bill.AccountB]++
	}

	sort.Slice(bills, func(i, j int) bool {
		return bills[i].Date.After(bills[j].Date)
	})

	b.data = bills
	b.countAmount()
}

func (b *Bill) countAmount() {
	if b.headButton == nil {
		return
	}

	countMap := make(map[string]float64)
	for _, v := range b.data {
		temp := map[string]float64{v.AccountA: v.AccountAPay, v.AccountB: v.AccountBPay}
		for k, vv := range temp {
			switch kk := strings.Split(k, ":")[0]; kk {
			case constant.Expenses:
				countMap[kk] = countMap[kk] + vv
			case constant.Liabilities, constant.Income:
				countMap[kk] = -countMap[kk] + vv
			}
		}
	}
	for i, v := range []string{constant.Income, constant.Expenses, constant.Liabilities} {
		text := b.headButton[i].Text
		text = strings.Split(text, ":")[0]
		text = fmt.Sprintf("%v:%s", text, strconv.FormatFloat(countMap[v], 'f', 2, 64))
		b.headButton[i].SetText(text)
		b.headButton[i].Refresh()
	}
}
