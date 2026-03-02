package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goggle-source/moneyLotServic/internal/config"
	"github.com/goggle-source/moneyLotServic/internal/lib/logger"
	"github.com/goggle-source/moneyLotServic/internal/models"
)

type DB struct {
	log *slog.Logger
	DB  *sql.DB
}

func InitDB(log *slog.Logger, cfg *config.Cfg) (*DB, error) {
	conn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.HostDB, cfg.DB.Port, cfg.DB.DbName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return &DB{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return &DB{}, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdle)
	db.SetConnMaxLifetime(cfg.DB.ConnLifeTime)

	return &DB{DB: db, log: log}, nil

}

func (d *DB) AddMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "postgresql.AddToMoneyUser"

	log := d.log.With(slog.String("op", op))

	log.Info("start AddMoneyToUser")

	log.Debug("data", slog.Float64("money", money), slog.String("userID", userID))

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error("error create transaction", logger.Err(err))
		return false, ValidationErrorPostgresql(err)
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
	INSERT INTO money (userID, userMoney) 
	VALUES ($1, $2)
	ON CONFLICT (userID) 
	DO UPDATE SET userMoney = money.userMoney + $2`,
		userID, money)

	if err != nil {
		log.Error("error update userMoney", logger.Err(err))
		return false, ValidationErrorPostgresql(err)
	}

	rowsAffected, err := res.RowsAffected()

	log.Debug("rows", slog.Int("r", int(rowsAffected)))

	if err != nil {
		log.Error("error get rowsAffected", logger.Err(err))
	}

	if rowsAffected == 0 {
		log.Error("error update userMoney", slog.Int("rows", int(rowsAffected)))
		return false, ValidationErrorPostgresql(err)
	}

	if err = tx.Commit(); err != nil {
		log.Error("error commit transaction", logger.Err(err))
		return false, ValidationErrorPostgresql(err)
	}

	log.Info("success AddMoneyToUser")

	return true, nil
}

func (d *DB) ReduceMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "postgresql.ReduceMoneyToUser"

	log := d.log.With(slog.String("op", op))

	log.Info("start ReduceMoneyToUser")

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error("error create transaction", logger.Err(err))
		return false, ValidationErrorPostgresql(err)
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
	UPDATE money SET userMoney = userMoney - $1
	WHERE userID = $2 AND userMoney >= $1`,
		money, userID)

	if err != nil {
		log.Error("error update userMoney", logger.Err(err))
		return false, ValidationErrorPostgresql(err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Error("error get rowsAffected", logger.Err(err))
		return false, ErrDatabase
	}

	if rowsAffected == 0 {
		log.Error("error update userMoney", slog.Int("rows", int(rowsAffected)))
		isCheckUser, err := d.CheckingIsUser(ctx, userID, tx)
		if err != nil {
			log.Error("error check user in database", logger.Err(err))
			return false, ValidationErrorPostgresql(err)
		}

		if !isCheckUser {
			log.Info("user is not found")
			return false, ErrNotFound
		}

		return false, ErrLackOfFunds
	}

	if err = tx.Commit(); err != nil {
		log.Error("error commit transaction", logger.Err(err))
		return false, ValidationErrorPostgresql(err)
	}

	log.Info("success ReduceMoneyToUser")
	return true, nil
}

func (d *DB) GetMoneyToUser(ctx context.Context, userID string) (float64, error) {
	const op = "postgresql.GetMoneyToUser"

	log := d.log.With(slog.String("op", op))

	log.Info("start GetMoneyToUser")

	var allMoney float64

	err := d.DB.QueryRowContext(ctx, `
	SELECT userMoney FROM money WHERE userID = $1`, userID).Scan(&allMoney)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		log.Error("error get userMoney", logger.Err(err))
		return 0, ValidationErrorPostgresql(err)
	}

	log.Info("success GetMoneyToUser")

	return allMoney, nil
}

func (d *DB) HealthCheack(ctx context.Context) (models.HealthCheakDatabase, error) {
	const op = "db.HealthCheack"

	log := d.log.With(slog.String("op", op))

	log.Info("start healthCheack")

	start := time.Now()

	_, err := d.DB.Exec("SELECT 1")
	if err != nil {
		log.Info("request execution error", logger.Err(err))
		return models.HealthCheakDatabase{ConnDB: false, LeadTimeDB: time.Since(start)}, err
	}
	log.Info("success request to DB")

	return models.HealthCheakDatabase{ConnDB: true, LeadTimeDB: time.Since(start)}, nil
}

func (d *DB) CheckingIsUser(ctx context.Context, userID string, tx *sql.Tx) (bool, error) {
	const op = "postgresql.CheckingIsUser"

	log := d.log.With(slog.String("op", op))

	log.Info("start checkingIsUser")

	var isValueUser bool

	err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM money WHERE userID = $1)", userID).Scan(&isValueUser)
	if err != nil {
		log.Error("error get userID from database", logger.Err(err))
		return false, err
	}
	log.Info("success checking an user")

	return isValueUser, nil
}
