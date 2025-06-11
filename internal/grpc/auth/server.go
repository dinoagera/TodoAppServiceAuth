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
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) (int64, error)
}
type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}
func (s *serverAPI) LoginUser(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}
	if len(req.GetPassword()) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password must not to be less than 6 symbol")
	}
	token, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error:%v", err)
	}
	return &sso.LoginResponse{
		Token:   token,
		Message: "User logined successufully",
	}, nil
}
func (s *serverAPI) RegisterUser(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}
	if len(req.GetPassword()) < 6 {
		return nil, status.Error(codes.InvalidArgument, "password must not to be less than 6 symbol")
	}
	uid, err := s.auth.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}
	return &sso.RegisterResponse{
		Message: fmt.Sprintf("User registered successufully. Your id:%d", uid),
	}, nil
}
