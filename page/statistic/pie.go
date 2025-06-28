package statistic

import (
	"bookkeeper/convert"
	"bookkeeper/service"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/wcharczuk/go-chart"
	"image"
	"log"
)

type Pie struct {
	content fyne.CanvasObject
	data    *service.Statistic
}

func NewPie() *Pie {
	return &Pie{
		data: service.StatisticService,
	}
}

func (p *Pie) Content() fyne.CanvasObject {
	if p.content != nil {
		return p.content
	}

	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: convert.MapToChartValues(p.data.ExpenseMap)}

	buffer := bytes.NewBuffer([]byte{})
	err := pie.Render(chart.PNG, buffer)
	if err != nil {
		log.Println(err)
	}

	decode, _, err := image.Decode(buffer)
	if err != nil {
		log.Println(err)
	}
	p.content = canvas.NewImageFromImage(decode)
	return p.content
}
