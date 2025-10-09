package ui

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var statementMap = map[fyne.CanvasObject]*statement{}
var statementId = map[int]*statement{}

type statement struct {
	date      *widget.Label
	usage     *widget.Label
	way       *widget.Label
	cost      *widget.Label
	update    *widget.Button
	delete    *widget.Button
	contianer fyne.CanvasObject
}

func createStatementList() *widget.List {
	c := widget.NewList(
		func() int { return len(service.FetchData()) },
		func() fyne.CanvasObject { return CreateStatement() },
		func(id widget.ListItemID, o fyne.CanvasObject) { UpdateStatement(id, o) },
	)

	c.OnSelected = func(id widget.ListItemID) {
		s := statementId[id]
		s.delete.OnTapped = func() {
			service.Delete(id)
		}
		s.delete.Show()
		s.update.Show()
	}
	c.OnUnselected = func(id widget.ListItemID) {
		c := statementId[id]
		c.delete.Hide()
		c.update.Hide()
	}

	return c
}

func UpdateStatement(i int, o fyne.CanvasObject) {
	statementMap[o].updateStatement(i)
	statementId[i] = statementMap[o]
}

func CreateStatement() fyne.CanvasObject {
	date := widget.NewLabel(constant.Zero)
	usage := widget.NewLabel(constant.Zero)
	way := widget.NewLabel(constant.Zero)
	cost := widget.NewLabel(constant.Zero)
	update := widget.NewButton(constant.Update, func() {})
	del := widget.NewButton(constant.Delete, func() {})
	del.Importance = widget.DangerImportance
	update.Importance = widget.HighImportance
	update.Hide()
	del.Hide()

	c := container.NewVBox(
		container.NewHBox(usage, date, layout.NewSpacer(), update, del),
		container.NewHBox(way, layout.NewSpacer(), cost),
	)

	statementMap[c] = &statement{date: date, usage: usage, way: way, cost: cost, update: update, delete: del, contianer: c}
	return c
}

func (s *statement) updateStatement(i int) {
	data := service.DataByIndex(i)
	s.date.SetText(data.Date.Format(time.DateOnly))
	s.usage.SetText(fmt.Sprintf("%s.%s", data.Terminal, data.Usage))
	s.way.SetText(fmt.Sprintf("%s---%s", data.From.Name, data.To.Name))
	s.cost.SetText(strconv.FormatFloat(data.To.Cost, 'f', 2, 64))
}
