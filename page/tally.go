package page

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/event"
	"bookkeeper/layoutCustom"
	"bookkeeper/service"
	"bookkeeper/ui"
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

func NewTally() service.Component {
	spend := widget.NewEntryWithData(service.TallyService.Cost)
	spend.TextStyle.Bold = true
	spend.SetPlaceHolder("输入花费金额 -- 9.9")

	date := widget.NewDateEntry()
	date.Bind(service.TallyService.Date)
	now := time.Now()
	date.SetDate(&now)

	payee := widget.NewEntryWithData(service.TallyService.Receiver)
	payee.SetPlaceHolder("商户/收款人")
	activity := widget.NewEntryWithData(service.TallyService.Usage)
	activity.SetPlaceHolder("用途 -- 午饭")
	u := container.New(layoutCustom.NewSplit(0.318), payee, activity)
	uu := &usage{concent: u, payee: payee, activity: activity}

	to := ui.AccountSelect(service.TallyService.To, constant.Liabilities)
	from := ui.AccountSelect(service.TallyService.From, constant.Expenses)

	finish := widget.NewButton("完成", func() {
		service.BillService.Add(service.TallyService.Deal())
		event.CurrentEvent = math.MaxInt32
	})
	finish.Importance = widget.HighImportance

	event.UiFuncMap[constant.UpdateEvent] = append(event.UiFuncMap[constant.UpdateEvent], func() {
		deal := service.BillService.Statements[<-event.DataIndex].Deals[<-event.DataIndex]
		spend.SetText(convert.Float64ToString(deal.Payment.Cost))
		date.SetDate(&deal.Date)
		uu.payee.SetText(deal.Payee)
		uu.activity.SetText(deal.Usage)

		service.TallyService.From.Prefix.Set(deal.Payment.Kind)
		service.TallyService.From.Account.Set(deal.Payment.Name)

		service.TallyService.To.Prefix.Set(deal.Receiver.Kind)
		service.TallyService.To.Account.Set(deal.Receiver.Name)

		service.BillService.Delete(deal)

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
