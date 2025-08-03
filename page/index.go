package page

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

type Index struct {
	bill    service.Component
	tally   service.Component
	account service.Component

	current int
}

func (i *Index) Content() fyne.CanvasObject {
	log.Println("start create index content")

	texts := []string{constant.Bill, constant.Tally, constant.Account}
	components := []service.Component{i.bill, i.tally, i.account}
	buttons := make([]*widget.Button, len(texts))
	var content *fyne.Container
	for ii, text := range texts {
		iii := ii
		buttons[ii] = widget.NewButton(text, func() {
			content.Remove(components[i.current].Content())
			buttons[i.current].Enable()
			i.current = iii
			content.Add(components[iii].Content())
			components[iii].Content().Refresh()
			buttons[iii].Disable()
		})
	}

	subContent := components[i.current].Content()
	bottom := container.NewGridWithColumns(len(components))
	for _, button := range buttons {
		bottom.Add(button)
	}
	content = container.NewBorder(nil, bottom, nil, nil, subContent)

	eventFunc[constant.Index] = func() {
		content.Remove(components[i.current].Content())
		buttons[i.current].Enable()
		i.current = 0
		content.Add(components[0].Content())
		components[0].Content().Refresh()
		buttons[0].Disable()
		content.Refresh()
		fmt.Println("refresh .......")
	}

	return content
}

func NewIndex() *Index {
	log.Println("create page index start...")
	return &Index{
		bill:    NewBill(),
		tally:   NewTally(),
		account: NewAccount(),
	}
}
