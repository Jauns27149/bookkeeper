package ui

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _aggregate = sync.OnceValue(func() *aggregate {
	aggregateService := service.GetAggregate()
	a := &aggregate{
		income:  widget.NewLabelWithData(aggregateService.Income),
		expense: widget.NewLabelWithData(aggregateService.Expenses),
		budget:  widget.NewLabelWithData(aggregateService.Budget),
	}
	a.income.Alignment, a.expense.Alignment, a.budget.Alignment =
		fyne.TextAlignCenter, fyne.TextAlignCenter, fyne.TextAlignCenter

	a._content = container.NewGridWithColumns(3,
		widget.NewLabelWithStyle(constant.IncomeChar, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle(constant.ExpensesChar, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle(constant.BudgetChar, fyne.TextAlignCenter, fyne.TextStyle{}),
		a.income, a.expense, a.budget)

	return a
})()

type aggregate struct {
	income   *widget.Label
	expense  *widget.Label
	budget   *widget.Label
	_content fyne.CanvasObject
}
