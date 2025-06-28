package statistic

import (
	"bookkeeper/intf"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Statistic struct {
	content fyne.CanvasObject
	pie     intf.Component
	items   intf.Component
}

func NewStatistic() *Statistic {
	return &Statistic{
		pie:   NewPie(),
		items: NewItems(),
	}
}

func (s *Statistic) Content() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.content = container.NewBorder(nil, s.items.Content(), nil, nil, s.pie.Content())
	return s.content
}
