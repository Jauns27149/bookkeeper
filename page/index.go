package page

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Index struct {
	bill    service.Component
	tally   service.Component
	account service.Component

	current int
}

func (i *Index) Content() fyne.CanvasObject {
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

	go i.listener(buttons)
	return content
}

func (i *Index) listener(buttons []*widget.Button) {
	for {
		select {
		case key := <-service.TallyService.Finish:
			var index int
			switch key {
			case constant.Bill:
				index = 0
			case constant.Tally:
				index = 1
			}
			fyne.Do(func() {
				buttons[index].OnTapped()
			})
		}
	}
}

func NewIndex() *Index {
	return &Index{
		bill:    NewBill(),
		tally:   NewTally(),
		account: NewAccount(),
	}
}
