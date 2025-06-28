package service

var (
	AccountService *Account
)

func Boot() {
	AccountService = NewAccount()
}
