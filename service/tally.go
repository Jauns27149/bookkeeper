package service

import (
	"bookkeeper/constant"
	"bookkeeper/convert"
	"bookkeeper/model"
	"log"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
)

var tally = &Tally{}

type Tally struct {
	pref    fyne.Preferences
	Account map[string][]string
}

func GetTally() *Tally {
	return tally
}

func Save(item model.Data) {
	tally.save(item)
}

func (t *Tally) save(item model.Data) {
	row := convert.DataToRow(item)
	key := item.Date.Format(constant.OnlyMonth)
	list := t.pref.StringList(key)
	list = append(list, row)
	t.pref.SetStringList(key, list)

	period := t.pref.StringList(constant.Period)
	if !slices.Contains(period, key) {
		t.pref.SetStringList(constant.Period, append(period, key))
	}
	log.Println("save data finished")
}

func (t *Tally) loadAccount() {
	accounts := t.pref.StringList(constant.Accounts)
	for _, a := range accounts {
		strs := strings.Split(a, ":")
		t.Account[strs[0]] = append(t.Account[strs[0]], strs[1])
	}

	log.Println("tally load account finished")
}

func(t *Tally) run() {
	t.pref=fyne.CurrentApp().Preferences()
	t.Account=make(map[string][]string)

	go tally.loadAccount()
}
