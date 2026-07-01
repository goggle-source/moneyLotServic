package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goggle-source/moneyLotServic/internal/config"
	"github.com/goggle-source/moneyLotServic/internal/domain"
	"github.com/goggle-source/moneyLotServic/internal/models"
)

type DB struct {
	DB *sql.DB
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

	return &DB{DB: db}, nil

}

func (d *DB) AddMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "postgresql.AddToMoneyUser"

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
	INSERT INTO money (userID, userMoney) 
	VALUES ($1, $2)
	ON CONFLICT (userID) 
	DO UPDATE SET userMoney = money.userMoney + $2`,
		userID, money)

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s:%w", op, domain.ErrAddMoney)
	}

	if err = tx.Commit(); err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return true, nil
}

func (d *DB) ReduceMoneyToUser(ctx context.Context, userID string, money float64) (bool, error) {
	const op = "postgresql.ReduceMoneyToUser"

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	defer tx.Rollback()
	var userMoney float64
	err = d.DB.QueryRowContext(ctx, "SELECT userMoney FROM money WHERE userID = $1", userID).Scan(&userMoney)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s:%w", op, domain.ErrNotFound)
		}
		return false, fmt.Errorf("%s:%w", op, err)
	}

	if userMoney < money {
		return false, fmt.Errorf("%s:%w", op, domain.ErrLackOfFunds)
	}

	res, err := tx.ExecContext(ctx, `
	UPDATE money SET userMoney = userMoney - $1
	WHERE userID = $2 AND userMoney >= $1`,
		money, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s:%w", op, domain.ErrNotFound)
		}
		return false, fmt.Errorf("%s:%w", op, err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	if rowsAffected == 0 {
		return false, fmt.Errorf("%s:%w", op, domain.ErrReduceMoney)
	}

	if err = tx.Commit(); err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return true, nil
}

func (d *DB) GetMoneyToUser(ctx context.Context, userID string) (float64, error) {
	const op = "postgresql.GetMoneyToUser"

	var allMoney float64

	err := d.DB.QueryRowContext(ctx, `
	SELECT userMoney FROM money WHERE userID = $1`, userID).Scan(&allMoney)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return allMoney, nil
}

func (d *DB) HealthCheack(ctx context.Context) (models.HealthCheakDatabase, error) {
	const op = "db.HealthCheack"

	start := time.Now()

	_, err := d.DB.Exec("SELECT 1")
	if err != nil {
		return models.HealthCheakDatabase{ConnDB: false, LeadTimeDB: time.Since(start)}, err
	}

	return models.HealthCheakDatabase{ConnDB: true, LeadTimeDB: time.Since(start)}, nil
}

func (d *DB) CheckingIsUser(ctx context.Context, userID string, tx *sql.Tx) (bool, error) {
	const op = "postgresql.CheckingIsUser"

	var isValueUser bool

	err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM money WHERE userID = $1)", userID).Scan(&isValueUser)
	if err != nil {
		return false, fmt.Errorf("%s:%w", op, err)
	}

	return isValueUser, nil
}
