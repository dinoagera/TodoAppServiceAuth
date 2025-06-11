package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"unicode/utf8"

	"github.com/dinoagera/api-auth/internal/domain/models"
	"github.com/dinoagera/api-auth/internal/lib/jwt"
	"github.com/dinoagera/api-auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
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
	if err := validateCredentials(email, password); err != nil {
		a.log.Warn("validation failed", "error", err)
		return "", fmt.Errorf("%w", ErrInvalidCredentials)
	}

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", "email", email)
			return "", fmt.Errorf("%w", ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", "email", email, "error", err)
		return "", fmt.Errorf("internal server error")
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid password", "email", email)
		return "", fmt.Errorf("%w", ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", "email", email, "error", err)
		return "", fmt.Errorf("failed to generate token")
	}

	a.log.Info("user logged in successfully", "email", email, "userID", user.ID)
	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	if err := validateCredentials(email, password); err != nil {
		a.log.Warn("validation failed", "error", err)
		return 0, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error("failed to generate password hash", "error", err)
		return 0, fmt.Errorf("failed to register user")
	}

	uid, err := a.usrSave.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user already exists", "email", email)
			return 0, fmt.Errorf("user with this email already exists")
		}
		a.log.Error("failed to save user", "email", email, "error", err)
		return 0, fmt.Errorf("failed to register user")
	}

	a.log.Info("user registered successfully", "email", email, "userID", uid)
	return uid, nil
}

func validateCredentials(email, password string) error {
	if email == "" {
		return ErrEmailRequired
	}
	if utf8.RuneCountInString(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}
