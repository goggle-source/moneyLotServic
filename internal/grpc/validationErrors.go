package grpc

import (
	"github.com/goggle-source/moneyLotServic/internal/business"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidationErrorsToBusiness(err error) error {
	arrErr := map[error]error{
		business.ErrDatabase:        status.Error(codes.Internal, "err database"),
		business.ErrConnDB:          status.Error(codes.Internal, "err database"),
		business.ErrBusiness:        status.Error(codes.Internal, "err application"),
		business.ErrInvalidArgument: status.Error(codes.InvalidArgument, "invalid field"),
		business.ErrLackOfFunds:     status.Error(codes.InvalidArgument, "lack of funds money"),
		business.ErrSmallDeposit:    status.Error(codes.InvalidArgument, "deposit don't small zero"),
	}

	result, ok := arrErr[err]
	if !ok {
		return status.Error(codes.Internal, "err application")
	}

	return result
}
