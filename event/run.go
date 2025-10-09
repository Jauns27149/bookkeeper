package event

import (
	"bookkeeper/constant"
	"log"

	"fyne.io/fyne/v2"
)

var eventFunc =map[uint]*event{
		constant.ConditionPrefixRefresh: &event{flag: make(chan struct{})},
		constant.ConditionSuffixRefresh: &event{flag: make(chan struct{})},
		constant.LoadBill:               &event{flag: make(chan struct{})},
		constant.BillRefresh:            &event{flag: make(chan struct{})},
	}

type event struct {
	flag chan struct{}
	fn   func()
}

func Run() {
	
}

func SetEventFunc(key uint, fn func()) {
	eventFunc[key].fn = fn
	close(eventFunc[key].flag)
}

func LaunchEvent(keys ...uint) {
	fyne.Do(func() {
		for _, key := range keys {
			if e, ok := eventFunc[key]; ok {
				<-e.flag
				e.fn()
			}
			log.Println("launch event func finished, key:", key)
		}
	})
}
