package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/barkimedes/go-deepcopy"
	"image/color"
	"log"
	"math"
	"time"
)

type Picker struct {
	padding  float32
	list     []string
	aligning chan float32
	current  int

	subSize fyne.Size
	scroll  *container.Scroll
}

func NewPicker(list []string) *Picker {
	return &Picker{
		list:     list,
		padding:  5,
		aligning: make(chan float32, 1),
	}
}
func (p *Picker) Content() *fyne.Container {
	box := container.New(layout.NewCustomPaddedVBoxLayout(p.padding))
	var width float32
	box.Add(widget.NewLabel(""))
	p.subSize = box.Objects[0].MinSize()
	for i, v := range p.list {
		label := widget.NewLabel(v)
		box.Add(label)
		if i == 0 {
			width = max(width, label.MinSize().Width)
		}
	}
	p.current = len(p.list)/2 + 1
	box.Add(widget.NewLabel(""))
	width *= 1.618
	scroll := container.NewVScroll(box)
	elementSize := box.Objects[0].MinSize()
	height := elementSize.Height*3 + 10
	size := fyne.NewSize(width, height)
	scroll.Resize(size)

	line := canvas.NewLine(color.White)
	line.StrokeWidth = 5
	line.Position1 = fyne.NewPos(0, elementSize.Height)
	line.Position2 = fyne.NewPos(width, elementSize.Height)
	lineHead := container.NewWithoutLayout(line)
	temp, err := deepcopy.Anything(line)
	if err != nil {
		log.Panicln(err)
	}
	lineTail := temp.(*canvas.Line)
	lineTail.Position1.Y = lineTail.Position1.Y + elementSize.Height + 10
	lineTail.Position2.Y = lineTail.Position2.Y + elementSize.Height + 10

	p.scroll = scroll
	//go p.align()
	scroll.OnScrolled = func(position fyne.Position) {
		p.aligning <- position.Y
	}

	box.Objects[p.current].(*widget.Label).Importance = widget.HighImportance
	scroll.ScrollToOffset(fyne.NewPos(0, float32(p.current-1)*(p.padding+p.subSize.Height)))
	return container.NewWithoutLayout(scroll, lineHead, container.NewWithoutLayout(lineTail))
}

func (p *Picker) align() {
	for {
		select {
		case height := <-p.aligning:
			//log.Println("height: ", height)
			select {
			case height = <-p.aligning:
				p.aligning <- height
			case <-time.After(time.Millisecond * 300):
				p.scrollToOffset(height)
			}
		}
	}
}

func (p *Picker) scrollToOffset(height float32) {
	fmt.Println("start height:", height)
	h := p.padding + p.subSize.Height
	fmt.Println("unit height:", h)
	var target float32
	for i, _ := range p.list {
		left := float32(i) * h

		right := float32(i+1) * h
		if i == len(p.list) {
			right = math.MaxFloat32
		}
		if height >= left && height <= right {
			if height-left > right-height {
				target = right
				i = i + 2
			} else {
				target = left
				i = i + 1
			}
			fyne.Do(func() {
				p.scroll.ScrollToOffset(fyne.NewPos(0, target))
				box := p.scroll.Content.(*fyne.Container)
				label := box.Objects[i].(*widget.Label)
				label.Importance = widget.HighImportance
				box.Objects[p.current].(*widget.Label).Importance = widget.MediumImportance
				p.current = i
				box.Refresh()
			})
			fmt.Println("target:", target)
			break
		}
	}
}
