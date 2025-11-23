package repository

import "errors"

var (
	ErrUserNotFound      = errors.New("resource not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
