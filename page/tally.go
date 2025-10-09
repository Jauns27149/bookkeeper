package page

import (
	"bookkeeper/constant-old"
	"bookkeeper/convert"
	"bookkeeper/event"
	"bookkeeper/layoutCustom"
	"bookkeeper/service-old"
	"bookkeeper/ui-old"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"math"
	"time"
)

type tally struct {
	spend    *widget.Entry
	date     *widget.DateEntry
	usage    *usage
	payment  fyne.CanvasObject
	receiver fyne.CanvasObject
	finish   *widget.Button
	content  fyne.CanvasObject
}
type usage struct {
	concent  *fyne.Container
	payee    *widget.Entry
	activity *widget.Entry
}

func (t *tally) Content() fyne.CanvasObject {
	if t.content != nil {
		return t.content
	}

	t.content = container.NewVBox(t.spend, t.usage.concent, t.payment, t.receiver, t.date, t.finish)
	return t.content
}

func NewTally() service_old.Component {
	spend := widget.NewEntryWithData(service_old.TallyService.Cost)
	spend.TextStyle.Bold = true
	spend.SetPlaceHolder("输入花费金额 -- 9.9")

	date := widget.NewDateEntry()
	date.Bind(service_old.TallyService.Date)
	now := time.Now()
	date.SetDate(&now)

	payee := widget.NewEntryWithData(service_old.TallyService.Receiver)
	payee.SetPlaceHolder("商户/收款人")
	activity := widget.NewEntryWithData(service_old.TallyService.Usage)
	activity.SetPlaceHolder("用途 -- 午饭")
	u := container.New(layoutCustom.NewSplit(0.318), payee, activity)
	uu := &usage{concent: u, payee: payee, activity: activity}

	to := ui.AccountSelect(service_old.TallyService.To, constant_old.Liabilities)
	from := ui.AccountSelect(service_old.TallyService.From, constant_old.Expenses)

	finish := widget.NewButton("完成", func() {
		service_old.BillService.Add(service_old.TallyService.Deal())
		event.CurrentEvent = math.MaxInt32
	})
	finish.Importance = widget.HighImportance

	event.UiFuncMap[constant_old.UpdateEvent] = append(event.UiFuncMap[constant_old.UpdateEvent], func() {
		deal := service_old.BillService.Statements[<-event.DataIndex].Deals[<-event.DataIndex]
		spend.SetText(convert.Float64ToString(deal.Payment.Cost))
		date.SetDate(&deal.Date)
		uu.payee.SetText(deal.Payee)
		uu.activity.SetText(deal.Usage)

		service_old.TallyService.From.Prefix.Set(deal.Payment.Kind)
		service_old.TallyService.From.Account.Set(deal.Payment.Name)

		service_old.TallyService.To.Prefix.Set(deal.Receiver.Kind)
		service_old.TallyService.To.Account.Set(deal.Receiver.Name)

		service_old.BillService.Delete(deal)

	})

	return &tally{
		spend:    spend,
		date:     date,
		usage:    uu,
		receiver: from,
		payment:  to,
		finish:   finish,
	}
}
