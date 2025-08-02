package layoutCustom

import "fyne.io/fyne/v2"

var _ fyne.Layout = (*split)(nil)

type split struct {
	offset float32
}

func (s split) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 2 {
		return
	}

	left := size.Width * s.offset
	right := size.Width - left
	height := min(objects[0].MinSize().Height, objects[1].MinSize().Height)
	objects[0].Resize(fyne.NewSize(left, height))
	objects[1].Resize(fyne.NewSize(right, height))

	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Move(fyne.NewPos(left, 0))
}

func (s split) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var height, width float32
	for _, o := range objects {
		if h := o.MinSize().Height; h > height {
			height = h
		}
		width += o.MinSize().Width
	}
	return fyne.NewSize(width, height)
}

func NewSplit(offset float32) fyne.Layout {
	return &split{offset: offset}
}
