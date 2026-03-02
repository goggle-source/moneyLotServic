package app

import (
	"log/slog"

	"github.com/goggle-source/moneyLotServic/cmd/migrate"
	"github.com/goggle-source/moneyLotServic/internal/app/grpcapp"
	"github.com/goggle-source/moneyLotServic/internal/business"
	"github.com/goggle-source/moneyLotServic/internal/config"
	"github.com/goggle-source/moneyLotServic/internal/database"
	"github.com/goggle-source/moneyLotServic/internal/lib/logger"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Cfg) *App {
	db, err := database.InitDB(log, cfg)
	if err != nil {
		panic(err)
	}

	err = migrate.RunMigrations(cfg.DB.User, cfg.DB.Password, cfg.DB.DbName, cfg.DB.Port, cfg.DB.HostDB)
	if err != nil {
		log.Info("error run migrations", logger.Err(err))
	}

	businessServic := business.Init(log, db)

	gRPCServic := grpcapp.InitGRPCApp(log, *businessServic, cfg.GRPC.Port)

	return &App{
		GRPCServer: gRPCServic,
	}
}
