package service

import (
	"bookkeeper/constant"
	"bookkeeper/model"
	"bookkeeper/util"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"sort"
	"strings"
)

type Tally struct {
	pref     fyne.Preferences
	accounts []model.SubAccount

	model.Record
}

func (t *Tally) Prefixes() (prefixes []string) {
	prefixes = make([]string, len(t.accounts))
	for i, v := range t.accounts {
		prefixes[i] = v.Prefix
	}
	return
}

func (t *Tally) Suffixes(prefix string) (suffixes []string) {
	for _, v := range t.accounts {
		if v.Prefix == prefix {
			suffixes = v.Suffixes
			break
		}
	}
	return
}

func (t *Tally) Deal() model.Deal {
	return model.Deal{
		Date:     util.BindTime(constant.FyneDate, t.Date),
		Payee:    util.BindSting(t.Receiver),
		Usage:    util.BindSting(t.Usage),
		Payment:  util.Account(t.From, t.Cost, true),
		Receiver: util.Account(t.To, t.Cost, false),
	}
}

func NewTally() *Tally {
	pref := fyne.CurrentApp().Preferences()
	list := pref.StringList(constant.Accounts)
	m := make(map[string][]string)
	for _, v := range list {
		values := strings.Split(v, constant.Colon)
		if len(values) != 2 {
			continue
		}
		if _, ok := m[values[0]]; !ok {
			m[values[0]] = []string{values[1]}
		} else {
			m[values[0]] = append(m[values[0]], values[1])
		}
	}

	accounts := make([]model.SubAccount, 0, len(m))
	for k, v := range m {
		accounts = append(accounts, model.SubAccount{Prefix: k, Suffixes: v})
	}

	sort.Slice(accounts, func(i, j int) bool { return len(accounts[i].Suffixes) > len(accounts[j].Suffixes) })

	return &Tally{
		pref:     pref,
		accounts: accounts,

		Record: model.Record{
			Cost:     binding.NewString(),
			Receiver: binding.NewString(),
			Usage:    binding.NewString(),
			From: model.BindAccount{
				Prefix:  binding.NewString(),
				Account: binding.NewString(),
			},
			To: model.BindAccount{
				Prefix:  binding.NewString(),
				Account: binding.NewString(),
			},
			Date: binding.NewString(),
		},
	}

}
