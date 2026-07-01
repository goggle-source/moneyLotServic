package domain

import "errors"

var (
	ErrLackOfFunds = errors.New("lack of funds")
	ErrNotFound    = errors.New("user is not found")
	ErrAddMoney    = errors.New("error add money")
	ErrSmalDeposit = errors.New("deposit must be greater than 0")
	ErrReduceMoney = errors.New("error reduce money")
)
