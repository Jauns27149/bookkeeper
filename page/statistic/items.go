package statistic

import (
	"bookkeeper/service"
	"bookkeeper/util"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"strconv"
)

type Items struct {
	content fyne.CanvasObject
	data    map[string]float64
}

func NewItems() *Items {
	return &Items{
		data: service.StatisticService.ExpenseMap,
	}
}

func (i *Items) Content() fyne.CanvasObject {
	if i.content != nil {
		return i.content
	}

	box := container.NewVBox()
	texts := make([]string, 0, len(i.data))
	for k, v := range i.data {
		text := fmt.Sprintf("%v:%v", k, strconv.FormatFloat(v, 'f', 2, 64))
		texts = append(texts, text)
	}

	util.SortObjectsBySum(texts)
	for _, text := range texts {
		box.Add(canvas.NewText(text, color.White))
	}
	i.content = box
	return i.content
}
