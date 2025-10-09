package page

import (
	"bookkeeper/constant-old"
	"bookkeeper/event"
	"bookkeeper/service-old"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

type Index struct {
	bill    service_old.Component
	tally   service_old.Component
	account service_old.Component

	current    int
	content    *fyne.Container
	components []service_old.Component
	buttons    []*widget.Button
}

func (i *Index) Content() fyne.CanvasObject {
	log.Println("start create index content")

	texts := []string{constant_old.Bill, constant_old.Tally, constant_old.Account}
	components := []service_old.Component{i.bill, i.tally, i.account}
	i.components = components

	buttons := make([]*widget.Button, len(texts))
	i.buttons = buttons

	var content *fyne.Container
	for ii, text := range texts {
		iii := ii
		buttons[ii] = widget.NewButton(text, func() {
			i.changePage(iii)
		})
	}

	subContent := components[i.current].Content()
	bottom := container.NewGridWithColumns(len(components))
	for _, button := range buttons {
		bottom.Add(button)
	}
	content = container.NewBorder(nil, bottom, nil, nil, subContent)

	service_old.PageEventFunc[constant_old.Index] = func() {
		i.changePage(0)
	}

	event.UiFuncMap[constant_old.UpdateEvent] = append(event.UiFuncMap[constant_old.UpdateEvent], func() {
		i.changePage(1)
	})

	i.content = content
	return content
}

func (i *Index) changePage(index int) {
	log.Println("current page index: ", i.current)

	i.content.Remove(i.components[i.current].Content())
	i.buttons[i.current].Enable()
	i.current = index
	i.content.Add(i.components[index].Content())
	i.buttons[index].Disable()
	i.content.Refresh()

	log.Printf("page change to index [%v] page", index)
}

func NewIndex() *Index {
	log.Println("create page index start...")
	return &Index{
		bill:    NewBill(),
		tally:   NewTally(),
		account: NewAccount(),
	}
}
