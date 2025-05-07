package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/dinoagera/api-auth/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	Port       string
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func New(log *slog.Logger, auth authgrpc.Auth, port string) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		Port:       port,
	}
}
func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.Port))
	if err != nil {
		a.log.Debug("didnt start server.", "port:", a.Port, "err:", err.Error())
		a.log.Info("didnt start server")
		return fmt.Errorf("%w", err)
	}
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%w", err)
	}
	a.log.Debug("starting server ", "Address:", l.Addr().String())
	a.log.Info("starting server")
	return nil
}
func (a *App) Stop() {
	a.log.Info("grpc server stopped")
	a.gRPCServer.GracefulStop()
}
