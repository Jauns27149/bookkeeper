package event

import (
    "fyne.io/fyne/v2"
    "log"
)



func uiEnventLinstner() {
    for {
        envent := <-UiEvent
        if fns, ok := UiFuncMap[envent]; ok {
            for _, fn := range fns {
                fyne.Do(func() {
                    fn()
                })
            }
        } else {
            log.Panic("not found func for envent ", envent)
        }
        log.Println("done ui enven: ", envent)
    }
}
