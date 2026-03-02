package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/goggle-source/moneyLotServic/internal/app"
	"github.com/goggle-source/moneyLotServic/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := InitLogger(cfg.Env)

	log.Info("start money servic", slog.Any("cfg", cfg))

	app := app.New(log, cfg)

	if err := app.GRPCServer.Run(); err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()

	log.Info("gracefully stopped")
}

func InitLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "debug":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
