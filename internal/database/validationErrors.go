package database

import (
	"database/sql"

	"github.com/lib/pq"
)

func ValidationErrorPostgresql(err error) error {
	arrErr := map[pq.ErrorCode]error{
		"23": ErrConstraintViolation,
		"08": ErrConnDB,
		"02": ErrNotFound,
	}

	errPq, ok := err.(*pq.Error)
	if !ok {
		err := ValidationErrorSql(err)
		return err
	}
	result, ok := arrErr[errPq.Code[0:2]]
	if !ok {
		return ErrDatabase
	}

	return result
}

func ValidationErrorSql(err error) error {
	arr := map[error]error{
		sql.ErrNoRows: ErrNotFound,
	}

	result, ok := arr[err]
	if !ok {
		return ErrDatabase
	}

	return result
}
