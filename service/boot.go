package service

var (
	AccountService *Account
	DataService    *Data
	TallyService   *Tally
)

func Boot() {
	DataService = NewData()
	DataService.Refresh()
	AccountService = NewAccount()
	TallyService = NewTally()

	listener()

}

func listener() {
	preference()
	account()
}
