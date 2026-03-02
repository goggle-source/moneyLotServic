package business

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/goggle-source/moneyLotServic/internal/lib/logger"
	"github.com/goggle-source/moneyLotServic/internal/models"
)

// нужно будет потом добавить во все методы проверку userID через метод authServic

type Database interface {
	AddMoneyToUser(ctx context.Context, userID string, money float64) (bool, error)
	ReduceMoneyToUser(ctx context.Context, userID string, money float64) (bool, error)
	GetMoneyToUser(ctx context.Context, userID string) (float64, error)
	HealthCheack(ctx context.Context) (models.HealthCheakDatabase, error)
}

type BusinessServic struct {
	log *slog.Logger
	DB  Database
}

func Init(log *slog.Logger, db Database) *BusinessServic {
	return &BusinessServic{
		log: log,
		DB:  db,
	}
}

func (b *BusinessServic) AddMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "business.AddMoneyToUser"

	log := b.log.With(slog.String("op", op))

	log.Info("start addMoneyToUser")

	if money < 0 {
		log.Error("error small deposit")
		return false, ErrSmallDeposit
	}

	ok, err := b.DB.AddMoneyToUser(ctx, userID, money)
	if err != nil {
		log.Error("error addMoney", logger.Err(err))
		return false, ValidationErrorsToRepositoryPostgresql(err)
	}
	log.Info("success AddMoneyToUser")

	return ok, nil
}

func (b *BusinessServic) ReduceMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "business.ReduceMoneyToUser"

	log := b.log.With(slog.String("op", op))

	log.Info("start ReduceMoneyToUser")

	if money < 0 {
		log.Error("error small deposit")
		return false, ErrSmallDeposit
	}

	ok, err := b.DB.ReduceMoneyToUser(ctx, userID, money)
	if err != nil {
		log.Error("error ReduceMoneyToUser", logger.Err(err))
		return false, ValidationErrorsToRepositoryPostgresql(err)
	}

	return ok, nil
}

func (b *BusinessServic) GetMoneyToUser(ctx context.Context, userID string) (float64, error) {
	const op = "business.GetMoneyToUser"

	log := b.log.With(slog.String("op", op))

	log.Info("start getMoneyToUser")

	AllMoney, err := b.DB.GetMoneyToUser(ctx, userID)
	if err != nil {
		log.Error("error GetMoneyToUser", logger.Err(err))
		return 0, ValidationErrorsToRepositoryPostgresql(err)
	}

	return AllMoney, nil
}

func (b *BusinessServic) Health(ctx context.Context) (map[string]any, error) {
	const op = "business.HealthCheack"

	log := b.log.With(slog.String("op", op))

	log.Info("start healthCheack")

	result := make(map[string]any, 5)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	result["alloc, MB"] = ConvertByteInMB(m.Alloc)
	result["allMemory requested from OS, MB"] = ConvertByteInMB(m.Sys)
	result["count NumGC"] = m.NumGC

	infoDB, err := b.DB.HealthCheack(ctx)
	if err != nil {
		log.Error("error get info DB", logger.Err(err))
		return result, ValidationErrorsToRepositoryPostgresql(err)
	}

	result["infoDB"] = infoDB

	log.Info("success healthCheack")

	return result, nil
}

func ConvertByteInMB(b uint64) uint64 {
	return b / 1024 / 1024
}
