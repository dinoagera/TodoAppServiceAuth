package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/dinoagera/api-auth/internal/app/grpc"
	"github.com/dinoagera/api-auth/internal/services/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort string, storagePath string, tokenTTL time.Duration) *App {
	authService := auth.New(log, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, _, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
