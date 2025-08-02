package page

import (
	"bookkeeper/constant"
	"bookkeeper/layoutCustom"
	"bookkeeper/model"
	"bookkeeper/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"strings"
	"time"
)

type tally struct {
	spend    *widget.Entry
	date     *widget.DateEntry
	usage    *fyne.Container
	payment  *fyne.Container
	receiver *fyne.Container
	finish   *widget.Button
	content  fyne.CanvasObject
}

func (t *tally) Content() fyne.CanvasObject {
	if t.content != nil {
		return t.content
	}

	t.content = container.NewVBox(t.spend, t.usage, t.payment, t.receiver, t.finish)
	return t.content
}

func NewTally() service.Component {
	server := service.TallyService

	spend := widget.NewEntry()
	spend.TextStyle.Bold = true
	spend.SetPlaceHolder("输入花费金额 -- 9.9")

	date := widget.NewDateEntry()
	now := time.Now()
	date.SetDate(&now)

	payee := widget.NewEntry()
	payee.SetPlaceHolder("商户/收款人")
	activity := widget.NewEntry()
	activity.SetPlaceHolder("用途 -- 午饭")
	usage := container.New(layoutCustom.NewSplit(0.318), payee, activity)

	prefixes := server.Prefixes()
	paymentPrefix := widget.NewSelect(prefixes, nil)
	paymentPrefix.SetSelectedIndex(0)
	paymentSuffix := widget.NewSelect(server.Suffixes(prefixes[0]), nil)
	paymentSuffix.SetSelectedIndex(0)
	payment := container.New(layoutCustom.NewSplit(0.318), paymentPrefix, paymentSuffix)

	receiverPrefix := widget.NewSelect(prefixes, nil)
	receiverPrefix.SetSelectedIndex(1)
	receiverSuffix := widget.NewSelect(server.Suffixes(prefixes[1]), nil)
	receiverSuffix.SetSelectedIndex(0)
	receiver := container.New(layoutCustom.NewSplit(0.318), receiverPrefix, receiverSuffix)

	save := func() {
		cost, err := strconv.ParseFloat(spend.Text, 64)
		if err != nil {
			log.Panicln(err)
		}
		paymentName := fmt.Sprintf("%s:%s", paymentPrefix.Selected, paymentSuffix.Selected)
		receiverName := fmt.Sprintf("%s:%s", receiverPrefix.Selected, receiverSuffix.Selected)
		deal := model.Deal{
			Date:     *date.Date,
			Payee:    payee.Text,
			Usage:    activity.Text,
			Payment:  model.Account{Name: paymentName, Cost: cost, Kind: "CNY"},
			Receiver: model.Account{Name: receiverName, Cost: -cost, Kind: "CNY"},
		}
		service.DataService.Add(deal)
	}
	finish := widget.NewButton("完成", save)

	finish.Importance = widget.HighImportance

	t := &tally{}

	go func() {
		for {
			select {
			case item := <-server.UpdateItem:
				date.SetDate(&item.Date)
				payee.SetText(item.Payee)
				activity.SetText(item.Usage)
				paymentPrefix.SetSelected(strings.Split(item.Payment.Name, constant.Colon)[0])
				paymentSuffix.SetSelected(strings.Split(item.Payment.Name, constant.Colon)[1])
				receiverPrefix.SetSelected(strings.Split(item.Receiver.Name, constant.Colon)[0])
				receiverSuffix.SetSelected(strings.Split(item.Receiver.Name, constant.Colon)[1])
				spend.SetText(strconv.FormatFloat(item.Payment.Cost, 'f', 2, 64))
				//t.content.Refresh()
			}
		}
	}()

	t.spend = spend
	t.date = date
	t.usage = usage
	t.payment = payment
	t.receiver = receiver
	t.finish = finish
	return t
}
