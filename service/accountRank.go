package service

type Account struct {
	accounts map[string]int
}

func (a Account) First() (first string) {
	amount, first := 0, ""
	for k, v := range a.accounts {
		if v >= amount {
			first = k
		}
	}
	return
}
