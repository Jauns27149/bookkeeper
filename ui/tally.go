package ui

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/event"
	"bookkeeper/model"
	"bookkeeper/service"
	"fmt"
	"log"
	"maps"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _tally = &tally{}

type tally struct {
	content fyne.CanvasObject

	spend  *widget.Entry
	date   *widget.DateEntry
	usage  *usage
	from   tallyAccount
	to     tallyAccount
	finish *widget.Button
}
type tallyAccount struct {
	prefix *widget.Select
	suffix *widget.Select

	content fyne.CanvasObject
}
type usage struct {
	concent  *fyne.Container
	payee    *widget.Entry
	activity *widget.Entry
}

func (t *tally) createContent() {
	if t.content != nil {
		return
	}

	t.content = container.NewVBox(t.spend, t.usage.concent, t.from.content, t.to.content, t.date, t.finish)
}

func (t *tally) createUI() {
	t.spend = widget.NewEntry()
	t.spend.TextStyle.Bold = true
	t.spend.SetPlaceHolder("输入花费金额")

	t.date = widget.NewDateEntry()
	now := time.Now()
	t.date.SetDate(&now)

	t.usage = &usage{
		payee:    widget.NewEntry(),
		activity: widget.NewEntry(),
	}
	t.usage.payee.SetPlaceHolder("商户/收款人")
	t.usage.activity.SetPlaceHolder("用途 -- 午饭")
	t.usage.concent = container.NewBorder(nil, nil, t.usage.payee, nil, t.usage.activity)

	t.createAccount()

	t.createFinish()

	log.Println("create tally ui finished")
}

func (t *tally) createAccount() {
	tallyService := service.GetTally()

	t.from = tallyAccount{
		prefix: widget.NewSelect([]string{}, nil),
		suffix: widget.NewSelect([]string{}, nil),
	}
	t.from.content = container.NewBorder(nil, nil, t.from.prefix, nil, t.from.suffix)

	t.to = tallyAccount{
		prefix: widget.NewSelect([]string{}, nil),
		suffix: widget.NewSelect([]string{}, nil),
	}
	t.to.content = container.NewBorder(nil, nil, t.to.prefix, nil, t.to.suffix)

	prefixes := slices.Sorted(maps.Keys(tallyService.Account))
	t.from.prefix.SetOptions(prefixes)
	t.from.prefix.SetSelectedIndex(2)
	t.to.prefix.SetOptions(prefixes)
	t.to.prefix.SetSelectedIndex(4)

	fn := func(s string, w *widget.Select) {
		w.SetOptions(tallyService.Account[s])
		w.SetSelectedIndex(0)
	}
	fn(t.from.prefix.Selected, t.from.suffix)
	fn(t.to.prefix.Selected, t.to.suffix)

	t.from.prefix.OnChanged = func(s string) { fn(s, t.from.suffix) }
	t.to.prefix.OnChanged = func(s string) { fn(s, t.to.suffix) }
}

func (t *tally) createFinish() {
	t.finish = widget.NewButton(constant.Finish, func() {
		data := model.Data{
			Date: *t.date.Date,
			From: model.Account{
				Name: fmt.Sprintf("%v:%v", t.from.prefix.Selected, t.from.suffix.Selected),
				Cost: -convert.StringToFloat64(t.spend.Text),
				Kind: constant.CNY,
			},
			To: model.Account{
				Name: fmt.Sprintf("%v:%v", t.to.prefix.Selected, t.to.suffix.Selected),
				Cost: convert.StringToFloat64(t.spend.Text),
				Kind: constant.CNY,
			},
			Terminal: t.usage.payee.Text,
			Usage:    t.usage.activity.Text,
		}
		service.Save(data)
		event.LaunchEvent(constant.LoadBill)

		_home.selectContent(0)
	})
	t.finish.Importance = widget.HighImportance
}

func init() {
	flag := make(chan struct{})
	go _tally.createUI()

	close(flag)
	setContent(constant.Tally, func() fyne.CanvasObject {
		<-flag
		_tally.createContent()
		return _tally.content
	})
}
