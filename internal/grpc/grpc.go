package grpc

import (
	"context"
	"log/slog"

	"github.com/goggle-source/MoneyLotProto/gen/go/money"
	"github.com/goggle-source/moneyLotServic/internal/lib/logger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	userID = "userID"
)

type BS interface {
	AddMoneyToUser(ctx context.Context, userID string, money float64) (bool, error)
	ReduceMoneyToUser(ctx context.Context, userID string, money float64) (bool, error)
	GetMoneyToUser(ctx context.Context, userID string) (float64, error)
	Health(ctx context.Context) (map[string]any, error)
}

type GrpcServic struct {
	money.UnimplementedMoneyServicServer
	BS  BS
	log *slog.Logger
}

func Register(grpcServer *grpc.Server, bs BS, log *slog.Logger) {
	money.RegisterMoneyServicServer(grpcServer, &GrpcServic{log: log, BS: bs})
}

func (g *GrpcServic) AddMoney(ctx context.Context, in *money.AddMoneyRequest) (*money.AddMoneyResponse, error) {
	const op = "grpc.AddMoney"
	log := g.log.With(slog.String("op", op))

	log.Info("start addMoney")

	log.Debug("data", slog.Float64("money", in.GetMoney()))

	ok, err := g.BS.AddMoneyToUser(ctx, ctx.Value(userID).(string), in.GetMoney())
	if err != nil {
		log.Error("error add money", logger.Err(err))
		return &money.AddMoneyResponse{}, ValidationErrorsToBusiness(err)
	}

	log.Info("success add money")
	return &money.AddMoneyResponse{
		Result: ok,
	}, nil
}

func (g *GrpcServic) ReduceMoney(ctx context.Context, in *money.ReduceMoneyRequest) (*money.ReduceMoneyResponse, error) {
	const op = "grpc.ReduceMoney"

	log := g.log.With(slog.String("op", op))

	log.Info("start reduceMoney")

	ok, err := g.BS.ReduceMoneyToUser(ctx, ctx.Value(userID).(string), float64(in.GetMoney()))
	if err != nil {
		log.Error("error reduceMoneyToUser", logger.Err(err))
		return &money.ReduceMoneyResponse{}, ValidationErrorsToBusiness(err)
	}

	return &money.ReduceMoneyResponse{
		Result: ok,
	}, nil
}

func (g *GrpcServic) GetMoneyUser(ctx context.Context, in *money.GetMoneyUserRequest) (*money.GetMoneyUserResponse, error) {
	const op = "grpc.GetMoney"

	log := g.log.With(slog.String("op", op))

	log.Info("start getMoneyUser")

	Allmoney, err := g.BS.GetMoneyToUser(ctx, ctx.Value(userID).(string))
	if err != nil {
		log.Error("error getMoneyUser", logger.Err(err))
		return &money.GetMoneyUserResponse{}, ValidationErrorsToBusiness(err)
	}

	log.Info("success getMoneyUser", slog.Float64("count", Allmoney))
	return &money.GetMoneyUserResponse{
		AllMoney: Allmoney,
	}, nil
}

func (g *GrpcServic) Health(ctx context.Context, in *money.HealthProductRequest) (*money.HealthProductResponse, error) {
	const op = "grpc.Health"

	log := g.log.With(slog.String("op", op))

	log.Info("start healthCheack")

	info, err := g.BS.Health(ctx)
	if err != nil {
		log.Error("error get info App", logger.Err(err))
		return &money.HealthProductResponse{}, ValidationErrorsToBusiness(err)
	}

	result := map[string]*anypb.Any{}

	for key, value := range info {
		result[key] = value.(*anypb.Any)
	}
	log.Info("success healthCheak")

	return &money.HealthProductResponse{Info: result}, nil
}
