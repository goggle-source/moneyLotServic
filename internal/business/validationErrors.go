package business

import "github.com/goggle-source/moneyLotServic/internal/database"

func ValidationErrorsToRepositoryPostgresql(err error) error {
	arrErr := map[error]error{
		database.ErrConnDB:              ErrConnDB,
		database.ErrConstraintViolation: ErrInvalidArgument,
		database.ErrDatabase:            ErrDatabase,
		database.ErrLackOfFunds:         ErrLackOfFunds,
		database.ErrInvalidArgument:     ErrInvalidArgument,
		database.ErrNotFound:            ErrInvalidArgument,
	}

	result, ok := arrErr[err]
	if !ok {
		return ErrBusiness
	}

	return result
}
