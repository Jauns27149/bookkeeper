package ui

import (
	"bookkeeper/constant"
	"bookkeeper/event"
	"bookkeeper/service"
	"maps"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

var _condition = &condition{}

type condition struct {
	content fyne.CanvasObject

	perfix *widget.Select
	suffix *widget.Select
	start  *widget.DateEntry
	end    *widget.DateEntry
	sure   *widget.Button
}

func getconditionContent() fyne.CanvasObject {
	_condition.createAccount()
	_condition.createtime()

	_condition.sure = widget.NewButton(constant.Sure, func() {
		event.LaunchEvent(constant.LoadBill)
	})
	_condition.sure.Importance = widget.HighImportance

	_condition.content = container.NewVBox(
		container.NewBorder(nil, nil, container.NewHBox(_condition.perfix, _condition.suffix), _condition.sure),
		container.NewGridWithColumns(2, _condition.start, _condition.end),
	)

	return _condition.content
}

func (c *condition) createtime() {
	condition := service.GetCondition()

	c.start = widget.NewDateEntry()
	c.start.OnChanged = func(t *time.Time) {
		service.GetCondition().Start = *t
	}
	c.start.SetDate(&condition.Start)

	c.end = widget.NewDateEntry()
	c.end.SetDate(&condition.End)
	c.end.OnChanged = func(t *time.Time) {
		service.GetCondition().End = *t
	}
}

func (c *condition) createAccount() {
	condition := service.GetCondition()
	c.perfix = widget.NewSelectWithData([]string{}, condition.Perfix)
	c.suffix = widget.NewSelectWithData([]string{}, condition.Perfix)
	c.perfix.Selected = ".*"
	c.suffix.Selected = ".*"

	event.SetEventFunc(constant.ConditionPrefixRefresh, func() {
		accounts := service.GetCondition().Account
		c.perfix.Options = slices.Sorted(maps.Keys(accounts))
	})
	event.SetEventFunc(constant.ConditionSuffixRefresh, func() {
		c.suffix.Options = service.GetCondition().Account[c.perfix.Selected]
	})

	condition.Perfix.AddListener(binding.NewDataListener(func() {
		event.LaunchEvent(constant.ConditionSuffixRefresh)
	}))
}
