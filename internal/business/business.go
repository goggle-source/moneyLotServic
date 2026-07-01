package business

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/goggle-source/moneyLotServic/internal/domain"
	"github.com/goggle-source/moneyLotServic/internal/models"
)

type Database interface {
	AddMoneyToUser(ctx context.Context, userID string, money float64) (bool, error)
	ReduceMoneyToUser(ctx context.Context, userID string, money float64) (bool, error)
	GetMoneyToUser(ctx context.Context, userID string) (float64, error)
	HealthCheack(ctx context.Context) (models.HealthCheakDatabase, error)
}

type BusinessServic struct {
	DB Database
}

func Init(log *slog.Logger, db Database) *BusinessServic {
	return &BusinessServic{
		DB: db,
	}
}

func (b *BusinessServic) AddMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "business.AddMoneyToUser"

	if money < 0 {
		return false, fmt.Errorf("it's money less 0: %s:%w", op, domain.ErrSmalDeposit)
	}

	ok, err := b.DB.AddMoneyToUser(ctx, userID, money)
	if err != nil {
		return false, fmt.Errorf("err add money in database layer: %s:%w", op, err)
	}

	return ok, nil
}

func (b *BusinessServic) ReduceMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "business.ReduceMoneyToUser"

	if money < 0 {
		return false, fmt.Errorf("it's money less 0: %s:%w", op, domain.ErrSmalDeposit)
	}

	ok, err := b.DB.ReduceMoneyToUser(ctx, userID, money)
	if err != nil {
		return false, fmt.Errorf("err reduce money in database layer: %s:%w", op, err)
	}

	return ok, nil
}

func (b *BusinessServic) GetMoneyToUser(ctx context.Context, userID string) (float64, error) {
	const op = "business.GetMoneyToUser"

	AllMoney, err := b.DB.GetMoneyToUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("err get money in database layer: %s:%w", op, err)
	}

	return AllMoney, nil
}

func (b *BusinessServic) Health(ctx context.Context) (map[string]any, error) {
	const op = "business.HealthCheack"

	result := make(map[string]any, 5)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	result["alloc, MB"] = ConvertByteInMB(m.Alloc)
	result["allMemory requested from OS, MB"] = ConvertByteInMB(m.Sys)
	result["count NumGC"] = m.NumGC

	infoDB, err := b.DB.HealthCheack(ctx)
	if err != nil {
		return result, fmt.Errorf("Err health check in database: %s:%w", op, err)
	}

	result["infoDB"] = infoDB

	return result, nil
}

func ConvertByteInMB(b uint64) uint64 {
	return b / 1024 / 1024
}
