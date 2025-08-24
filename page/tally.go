package page

import (
	"bookkeeper/constant"
	"bookkeeper/layoutCustom"
	"bookkeeper/service"
	"bookkeeper/ui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

type tally struct {
	spend    *widget.Entry
	date     *widget.DateEntry
	usage    *fyne.Container
	payment  fyne.CanvasObject
	receiver fyne.CanvasObject
	finish   *widget.Button
	content  fyne.CanvasObject
}

func (t *tally) Content() fyne.CanvasObject {
	if t.content != nil {
		return t.content
	}

	t.content = container.NewVBox(t.spend, t.usage, t.payment, t.receiver, t.date, t.finish)
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
	usage := container.New(layoutCustom.NewSplit(0.318), payee, activity)

	to := ui.AccountSelect(service.TallyService.To, constant.Liabilities)
	from := ui.AccountSelect(service.TallyService.From, constant.Expenses)

	finish := widget.NewButton("完成", func() {
		service.BillService.Add(service.TallyService.Deal())
	})
	finish.Importance = widget.HighImportance

	return &tally{
		spend:    spend,
		date:     date,
		usage:    usage,
		receiver: from,
		payment:  to,
		finish:   finish,
	}
}
