package service

var (
	AccountService *Account
	DataService    *Data
)

func Boot() {
	DataService = NewData()

	DataService.Refresh()

	AccountService = NewAccount()

	listener()

}

func listener() {
	preference()
	account()
}
