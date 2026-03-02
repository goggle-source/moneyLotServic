package business

import "errors"

var (
	ErrConnDB          = errors.New("couldn't connect to the database")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrDatabase        = errors.New("error database")
	ErrBusiness        = errors.New("error app")
	ErrLackOfFunds     = errors.New("lack of funds userMoney")
	ErrSmallDeposit    = errors.New("The deposit must be greater than 0")
)
