package event

var UiEvent = make(chan uint, 1)
var UiFuncMap = make(map[uint][]func())
var DataIndex = make(chan int, 2)
var CurrentEvent uint = 0

func Run() {
    go uiEnventLinstner()
}
