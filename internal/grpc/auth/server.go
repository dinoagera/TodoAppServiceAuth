package authgrpc

import (
	"context"
	"fmt"

	sso "github.com/dinoagera/proto/gen/go/myservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (int64, error)
}
type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}
func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is empty")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}
	if len(req.GetPassword()) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password must not to be less than 6 symbol")
	}
	token, err := s.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &sso.LoginResponse{
		Token:   token,
		Message: "User logined successufully",
	}, nil
}
func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "username is empty")
	}
	if len(req.GetUsername()) < 3 {
		return nil, status.Error(codes.InvalidArgument, "your name must not to be less than 3 symbol")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}
	if len(req.GetPassword()) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password must not to be less than 6 symbol")
	}
	uid, err := s.auth.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &sso.RegisterResponse{
		Message: fmt.Sprintf("User registered successufully. Your id:%d", uid),
	}, nil
}
