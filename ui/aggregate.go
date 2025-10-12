package ui

import (
	"bookkeeper/constant"
	"bookkeeper/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _aggregate = &aggregate{}

type aggregate struct {
	income   *widget.Label
	expense  *widget.Label
	budget   *widget.Label
	content fyne.CanvasObject
}

func(a *aggregate) run(){
	aggregateService := service.GetAggregate()
		a.income=  widget.NewLabelWithData(aggregateService.Income)
		a.expense=widget.NewLabelWithData(aggregateService.Expenses)
		a.budget=  widget.NewLabelWithData(aggregateService.Budget)
	
	a.income.Alignment, a.expense.Alignment, a.budget.Alignment =
		fyne.TextAlignCenter, fyne.TextAlignCenter, fyne.TextAlignCenter

	a.content = container.NewGridWithColumns(3,
		widget.NewLabelWithStyle(constant.IncomeChar, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle(constant.ExpensesChar, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle(constant.BudgetChar, fyne.TextAlignCenter, fyne.TextStyle{}),
		a.income, a.expense, a.budget)
}
