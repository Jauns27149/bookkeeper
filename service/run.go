package service

import "log/slog"

func Run() {
	_bill.run()
	accounts.run()
	tally.run()
}

func AddAccount(account string) {
	accounts.save(account)
	slog.Info("save account success.", "account", account)
}
