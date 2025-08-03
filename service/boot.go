package service

import "log"

var (
	AccountService *Account
	BillService    *Bill
	TallyService   *Tally
)

func Boot() {
	BillService = NewBill()
	go BillService.Load()

	AccountService = NewAccount()
	TallyService = NewTally()

	listener()
	log.Println("data servicer start successful")
}

func listener() {
	go dataEvent()
}
