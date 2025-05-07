package app

import (
	"log/slog"
	"os"
	"time"

	grpcapp "github.com/dinoagera/api-auth/internal/app/grpc"
	"github.com/dinoagera/api-auth/internal/services/auth"
	"github.com/dinoagera/api-auth/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort string, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.New(storagePath)
	if err != nil {
		log.Debug("failed to init db,", "storagePath:", storagePath)
		os.Exit(1)
	}
	authService := auth.New(log, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
