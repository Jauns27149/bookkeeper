package page

import (
	"bookkeeper/constant"
	"bookkeeper/service"
	"fyne.io/fyne/v2"
)

var eventFunc = make(map[uint]func())

func uiEvent() {
	for {
		switch <-service.BillService.UiEvent {
		case constant.Index:
			fyne.Do(func() {
				if fn := eventFunc[constant.Index]; fn != nil {
					fn()
				}
			})
		}
	}
}
