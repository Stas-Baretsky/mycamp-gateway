package auth

import (
	"context"
	"fmt"
	"gateway-api/internal/config"
	"gateway-api/internal/handlers/url/register"
	"gateway-api/internal/lib/logger/sl"
	"log/slog"
	"net"
	"strconv"
	"time"

	ssov1 "github.com/Stas-Baretsky/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	api ssov1.AuthClient
	log *slog.Logger
}

const grpcHost = "localhost"

// /TODO:добавить в конфиг таймаут grpc запроса и исправить инициализацию контекста
func New(log *slog.Logger, cfg config.Config) (*AuthClient, error) {
	const op = "client.auth.New"

	сс, err := grpc.NewClient(
		grpcAddress(&cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &AuthClient{
		api: ssov1.NewAuthClient(сс),
		log: log,
	}, nil
}

func (auth *AuthClient) Register(
	ctx context.Context,
	params register.RegisterParams,
) (int64, error) {
	const op = "client.auth.Register"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("email", params.Email),
	)

	ctx, _ = context.WithTimeout(ctx, time.Duration(5*time.Second))

	log.Info("enter register grpc call")

	respReg, err := auth.api.Register(ctx, &ssov1.RegisterRequest{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		auth.log.Warn("register failed", sl.Err(err))
	}
	log.Info("register ok")
	return respReg.UserId, nil
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.RPCAuthClient.Port))
}
