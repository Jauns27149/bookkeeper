package page

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Bottom struct {
	Items   []string
	Buttons []*widget.Button
}

func (b *Bottom) Content() fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(b.Items))
	for i, button := range b.Buttons {
		objects[i] = button
	}
	return container.NewGridWithColumns(len(objects), objects...)
}

func NewBottom(items []string) *Bottom {
	buttons := make([]*widget.Button, len(items))
	for i, item := range items {
		buttons[i] = widget.NewButton(item, func() {})
	}
	return &Bottom{
		Items:   items,
		Buttons: buttons,
	}
}
