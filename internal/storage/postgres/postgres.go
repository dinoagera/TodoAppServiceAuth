package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dinoagera/api-auth/internal/domain/models"
	"github.com/dinoagera/api-auth/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	db *pgx.Conn
}

func New(storagePath string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed connect to db, err:%w", err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed connect to db, err:%w", err)
	}
	return &Storage{db: conn}, nil
}
func (s *Storage) SaveUser(ctx context.Context, email string, PassHash []byte) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx, "INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id", email, PassHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%w", storage.ErrUserExists)
		}
		return 0, fmt.Errorf("failed to add user, err:%w", err)
	}
	return id, nil
}
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		ctx,
		"SELECT id, email, pass_hash FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
