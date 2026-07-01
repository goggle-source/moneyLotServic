package grpc

import (
	"errors"

	"github.com/goggle-source/moneyLotServic/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidationErrorsToBusiness(err error) error {
	if errors.Is(err, domain.ErrLackOfFunds) {
		return status.Error(codes.InvalidArgument, "lack of funds money")
	}
	if errors.Is(err, domain.ErrSmalDeposit) {
		return status.Error(codes.InvalidArgument, "deposit don't small zero")
	}
	if errors.Is(err, domain.ErrAddMoney) {
		return status.Error(codes.InvalidArgument, "error add money")
	}
	if errors.Is(err, domain.ErrNotFound) {
		return status.Error(codes.NotFound, "user not found")
	}
	return err
}
