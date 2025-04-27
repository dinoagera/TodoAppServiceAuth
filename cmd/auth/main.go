package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dinoagera/api-auth/internal/app"
	"github.com/dinoagera/api-auth/internal/config"
	"github.com/dinoagera/api-auth/internal/logger"
)

func main() {
	logger := logger.InitLogger()
	cfg := config.InitConfig(logger)
	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCSrv.Stop()
	logger.Info("application is stopped")
}
