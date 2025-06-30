package component

import (
	"bookkeeper/constant"
	"bookkeeper/model"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"time"
)

var updateItem = make(chan model.Deal, 1)

type Record struct {
	button  *widget.Button
	add     fyne.CanvasObject
	data    *service.Data
	content fyne.CanvasObject
}

// NewRecord TODO here is too long, need to split
func NewRecord() *Record {
	data := service.DataService
	button := widget.NewButton(constant.Tally, nil)
	box := container.NewVBox()

	dateEntry := widget.NewDateEntry()
	now := time.Now()
	dateEntry.SetDate(&now)
	box.Add(dateEntry)

	payees := data.Payees
	payee := widget.NewSelectEntry(payees)
	payee.SetPlaceHolder("收款人/商户/收入来源渠道")
	if len(payees) > 0 {
		payee.SetText(payees[0])
	}
	box.Add(payee)

	detail := widget.NewEntry()
	detail.SetPlaceHolder("详细描述，记录细节")
	box.Add(detail)

	accounts, err := data.Accounts.Get()
	if err != nil {
		log.Println(err)
	}

	accountSelect := make([]*widget.Select, 2)
	amounts := make([]*widget.Entry, 2)
	for i := range 2 {
		accountSelect[i] = widget.NewSelect(accounts, nil)
		account := accountSelect[i]
		if accounts != nil && len(accounts) > i {
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

			costA, _ := strconv.ParseFloat(amounts[0].Text, 64)
			costB, _ := strconv.ParseFloat(amounts[1].Text, 64)
			deal := model.Deal{
				Date:         *dateEntry.Date,
				Payee:        payee.Text,
				Usage:        detail.Text,
				AccountA:     accountSelect[0].Selected,
				AccountAPay:  costA,
				AccountAKind: "CNY",
				AccountB:     accountSelect[1].Selected,
				AccountBPay:  costB,
				AccountBKind: "CNY",
			}

			amounts[0].SetText("")
			amounts[1].SetText("")
			detail.Text = ""

			box.Hide()
			data.GetDeals("")
			data.Deals = append(data.Deals, deal)
			data.Save(deal)
			data.RefreshPage()
		}
		acc := container.NewGridWithColumns(2)
		acc.Add(account)
		acc.Add(amount)
		box.Add(acc)
	}

	data.Accounts.AddListener(binding.NewDataListener(func() {
		list, err := data.Accounts.Get()
		if err != nil {
			log.Println(err)
		}

		for i, v := range accountSelect {
			v.SetOptions(list)
			if len(list) > i {
				v.SetSelectedIndex(i)
			}
		}
	}))

	// update item
	go func() {
		for {
			item := <-updateItem
			fyne.Do(func() {
				box.Show()
				dateEntry.SetDate(&item.Date)
				payee.SetText(item.Payee)
				detail.SetText(item.Usage)
				accountSelect[0].SetSelected(item.AccountA)
				amounts[0].SetText(strconv.FormatFloat(item.AccountAPay, 'f', 2, 64))
				accountSelect[1].SetSelected(item.AccountB)
				amounts[1].SetText(strconv.FormatFloat(item.AccountBPay, 'f', 2, 64))
				data.RefreshPage()
			})
		}
	}()

	return &Record{
		button: button,
		data:   data,
		add:    box,
	}
}

func (r *Record) Content() fyne.CanvasObject {
	if r.content != nil {
		return r.content
	}

	r.add.Hide()
	r.button.OnTapped = func() {
		if r.add.Visible() {
			r.add.Hide()
		} else {
			r.add.Show()
		}
	}

	r.content = container.NewVBox(r.button, r.add)
	return r.content
}
