package bill

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

type Gather struct {
	content   fyne.CanvasObject
	income    *canvas.Text
	expense   *canvas.Text
	liability *canvas.Text
	data      *service.Data
}

func NewGather() *Gather {
	g := &Gather{
		data: service.DataService,
	}

	value, _ := g.data.Income.Get()
	g.income = canvas.NewText(fmt.Sprintf("%v: %v", constant.Income, value), constant.Red)
	g.income.Alignment = fyne.TextAlignCenter
	value, _ = g.data.Expense.Get()
	g.expense = canvas.NewText(fmt.Sprintf("%v: %v", constant.Expenses, value), constant.Green)
	g.expense.Alignment = fyne.TextAlignCenter
	value, _ = g.data.Liability.Get()
	g.liability = canvas.NewText(fmt.Sprintf("%v: %v", constant.Liabilities, value), constant.Green)
	g.liability.Alignment = fyne.TextAlignCenter

	g.data.Income.AddListener(binding.NewDataListener(func() {
		value, _ = g.data.Income.Get()
		g.income.Text = fmt.Sprintf("%v: %v", constant.Income, value)
	}))
	g.data.Expense.AddListener(binding.NewDataListener(func() {
		value, _ = g.data.Expense.Get()
		g.expense.Text = fmt.Sprintf("%v: %v", constant.Expenses, value)
	}))
	g.data.Liability.AddListener(binding.NewDataListener(func() {
		value, _ = g.data.Liability.Get()
		g.liability.Text = fmt.Sprintf("%v: %v", constant.Liabilities, value)
	}))

	return g
}

func (g *Gather) Content() fyne.CanvasObject {
	if g.content != nil {
		return g.content
	}

	g.content = container.NewGridWithColumns(3, g.income, g.expense, g.liability)
	return g.content
}
