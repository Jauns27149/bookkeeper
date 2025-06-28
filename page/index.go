package page

import (
	"bookkeeper/constant"
	"bookkeeper/page/bill"
	"bookkeeper/page/statistic"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type Index struct {
	bottom     *Bottom
	components []fyne.CanvasObject
	current    int
}

func (i *Index) Content() fyne.CanvasObject {
	var content *fyne.Container

	contentMap := make(map[string]fyne.CanvasObject, len(i.bottom.Items))
	for ii := range len(contentMap) {
		contentMap[i.bottom.Items[ii]] = i.components[ii]
	}
	for ii, button := range i.bottom.Buttons {
		iii := ii
		button.OnTapped = func() {
			content.Remove(i.components[i.current])
			content.Add(i.components[iii])
			i.bottom.Buttons[iii].Disable()
			i.bottom.Buttons[i.current].Enable()
			i.current = iii
		}
	}

	i.bottom.Buttons[i.current].Disable()
	bottom := i.bottom.Content()
	object := i.components[i.current]
	content = container.NewBorder(nil, bottom, nil, nil, object)
	return content
}

func NewIndex() *Index {
	texts := []string{
		constant.Bill,
		constant.Statistic,
		constant.Account,
	}

	return &Index{
		bottom: NewBottom(texts),
		components: []fyne.CanvasObject{
			bill.NewBill().Content(),
			statistic.NewStatistic().Content(),
			NewAccount().Content(),
		},
	}
}
