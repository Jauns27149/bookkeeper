package base

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ButtonListH struct {
	buttons []*widget.Button
}

func (s *ButtonListH) Content() fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, len(s.buttons))
	for i, b := range s.buttons {
		objects[i] = b
	}
	return container.NewHBox(objects...)
}

func NewButtonListH(texts []string) *ButtonListH {
	buttons := make([]*widget.Button, len(texts))
	for i, text := range texts {
		buttons[i] = widget.NewButton(text, nil)
	}

	return &ButtonListH{
		buttons: buttons,
	}
}
