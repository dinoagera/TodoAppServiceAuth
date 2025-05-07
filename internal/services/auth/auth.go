package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dinoagera/api-auth/internal/domain/models"
	"github.com/dinoagera/api-auth/internal/lib/jwt"
	"github.com/dinoagera/api-auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log         *slog.Logger
	usrSave     UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
}
type UserSaver interface {
	SaveUser(ctx context.Context, email string, PassHash []byte) (int64, error)
}
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

func New(log *slog.Logger, usrSave UserSaver, usrProvider UserProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		usrSave:     usrSave,
		usrProvider: usrProvider,
		tokenTTL:    tokenTTL,
	}
}
func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", "err:", err.Error())
			return "", fmt.Errorf("%w", ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", "err:", err.Error())
		return "", fmt.Errorf("%w", ErrInvalidCredentials)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", "err:", err.Error())
		return "", fmt.Errorf("%w", ErrInvalidCredentials)
	}
	a.log.Info("User is login successufully")
	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token")
		return "", fmt.Errorf("%w", err)
	}
	return token, nil
}
func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Debug("failed to generate hash password")
		return 0, fmt.Errorf("%w", err)
	}
	uid, err := a.usrSave.SaveUser(ctx, email, passHash)
	if err != nil {
		a.log.Info("User is not registered")
		return 0, fmt.Errorf("%w", err)
	}
	a.log.Info("User is register successufully")
	return uid, nil
}
