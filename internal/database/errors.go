package database

import "errors"

var (
	ErrDatabase            = errors.New("error database")
	ErrConnDB              = errors.New("couldn't connect to the database")
	ErrNotFound            = errors.New("user is not found")
	ErrConstraintViolation = errors.New("violation of the restriction")
	ErrLackOfFunds         = errors.New("lack of funds")
	ErrInvalidArgument     = errors.New("invalid argument")
)
