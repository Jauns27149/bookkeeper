package tally

import (
	"bookkeeper/service"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type tally struct {
	usage   service.Component
	payment *fyne.Container
	receipt *fyne.Container
	cost    *fyne.Container
}

func (t *tally) Content() fyne.CanvasObject {
	return widget.NewPopUp(widget.NewLabel("here is tally"), fyne.CurrentApp().Driver().AllWindows()[0].Canvas())

	//v := reflect.ValueOf(t)
	//if v.Kind() == reflect.Ptr {
	//	v = v.Elem()
	//}
	//var objects []fyne.CanvasObject
	//for i := 0; i < v.NumField(); i++ {
	//	f := v.Field(i)
	//	if c, ok := f.Interface().(*fyne.Container); !f.IsNil() && ok {
	//		objects = append(objects, c)
	//	}
	//}
	//return container.NewVBox(objects...)
}

func NewTally() service.Component {
	return &tally{
		//usage: NewUsage(),
	}
}
