package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/goggle-source/moneyLotServic/internal/business"
	grpcServ "github.com/goggle-source/moneyLotServic/internal/grpc"
	"github.com/goggle-source/moneyLotServic/internal/grpc/interceptor"
	"github.com/goggle-source/moneyLotServic/internal/lib/logger"
	"google.golang.org/grpc"
)

type App struct {
	log  *slog.Logger
	GRPC *grpc.Server
	Port int
}

func InitGRPCApp(log *slog.Logger, bs business.BusinessServic, port int) *App {
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.LoggingInterceptor,
		interceptor.AuthInterceptor,
	))

	grpcServ.Register(gRPCServer, &bs, log)

	return &App{
		log:  log,
		GRPC: gRPCServer,
		Port: port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(slog.String("op", op))

	log.Info("start run server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.Port))
	if err != nil {
		log.Error("error listen TCP port", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("start grpc server", slog.String("port", l.Addr().String()))

	if err := a.GRPC.Serve(l); err != nil {
		log.Error("error start grpc server", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(slog.String("op", op))

	log.Info("stopping grpc server", slog.Int("port", a.Port))

	a.GRPC.GracefulStop()
}
