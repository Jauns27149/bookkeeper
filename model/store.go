package model

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ListStore struct {
	interfaceObjectMap map[fyne.CanvasObject]any
	idInterfaceMap     map[widget.ListItemID]fyne.CanvasObject
}

func NewListStore() *ListStore {
	return &ListStore{
		interfaceObjectMap: make(map[fyne.CanvasObject]any),
		idInterfaceMap:     make(map[widget.ListItemID]fyne.CanvasObject),
	}
}

func (l *ListStore) Store(key any, value any) {
	if k, ok := key.(fyne.CanvasObject); ok {
		l.interfaceObjectMap[k] = value
	}

	if k, ok := key.(widget.ListItemID); ok {
		l.idInterfaceMap[k] = value.(fyne.CanvasObject)
		return
	}
}

func (l *ListStore) Load(key any) any {
	if k, ok := key.(widget.ListItemID); ok {
		return l.interfaceObjectMap[l.idInterfaceMap[k]]
	}

	if k, ok := key.(fyne.CanvasObject); ok {
		return l.interfaceObjectMap[k]
	}

	return nil
}
