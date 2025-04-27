package storage

import "errors"

var (
	ErrUserExists   = errors.New("user aldready exists")
	ErrUserNotFound = errors.New("user not found")
)
