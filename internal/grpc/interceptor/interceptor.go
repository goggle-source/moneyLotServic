package interceptor

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/goggle-source/moneyLotServic/internal/config"
	publickey "github.com/goggle-source/moneyLotServic/internal/lib/PublicKey"
	"github.com/goggle-source/moneyLotServic/internal/lib/logger"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	userID           = "userID"
	ClientServicName = "client-servic-name"
)

var (
	publicRSKey   *rsa.PublicKey
	secretKeyOnce sync.Once
)

type Claims struct {
	Userid string
	jwt.RegisteredClaims
}

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// нужно будет потом будет сделать логгер так, чтобы не создавался новый логгер при каждом запросе
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	log.Info("request processing starts", slog.String("path", info.FullMethod))

	resp, err = handler(ctx, req)
	if err != nil {
		log.Error("request processing error", logger.Err(err))
	}

	log.Info("end of request processing", slog.Any("result", resp))

	return resp, err
}

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	name := md.Get(ClientServicName)

	if CheckingServicForAccessWithoutJWT(name[0]) {

		return handler(ctx, req)
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		GetSecretKey()
		return publicRSKey, nil

	}, jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithIssuedAt())
	// потом нужно будет добавить проверку на issuer(тот, выдал токен)

	if err != nil || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	ctx = context.WithValue(ctx, userID, claims.Userid)

	return handler(ctx, req)
}

func CheckingServicForAccessWithoutJWT(name string) bool {
	arr := map[string]bool{
		"productLotServic": true,
		"auctionLotServic": true,
	}

	return arr[name]
}

func GetSecretKey() {
	secretKeyOnce.Do(func() {
		key, err := publickey.LoadPublicKey(config.MustLoad().Path)
		if err != nil {
			fmt.Println(err)
		}
		publicRSKey = key
	})

}
