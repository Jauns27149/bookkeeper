package ui

import (
	"bookkeeper/constant"
	"log"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _home =&home{}

type home struct {
	objects []*object
	content *fyne.Container
	current int
}

type object struct {
	button  *widget.Button
	content func() fyne.CanvasObject
}

func (h *home) selectContent(i int) {
	h.content.Remove(h.objects[h.current].content())
	h.objects[h.current].button.Importance = widget.MediumImportance
	h.current = i
	h.content.Add(h.objects[h.current].content())
	h.objects[h.current].button.Importance = widget.HighImportance
	h.content.Refresh()
}

func(h *home) run(){
	for i, v := range []string{constant.Bill, constant.Tally, constant.Account} {
		h.objects = append(h.objects, &object{
			button: widget.NewButton(v, func() {
				h.selectContent(i)
			}),
		})
	}

	_bill.run()
	_tally.run()
	_accounts.run()
	slog.Info("home ui init finished")
}

func Content() fyne.CanvasObject {
	if _home.content == nil {
		c := container.NewGridWithColumns(len(_home.objects))
		for _, o := range _home.objects {
			c.Add(o.button)
		}
		_home.content = container.NewBorder(nil, c, nil, nil, _home.objects[_home.current].content())
		_home.objects[_home.current].button.Importance = widget.HighImportance
	}

	log.Println("home fetch init content success")
	return _home.content
}

func setContent(name string, content func() fyne.CanvasObject) {
	for _, o := range _home.objects {
		if o.button.Text == name {
			o.content = content
			log.Println("home set object content finished, name:", name)
			break
		}
	}
}
