package repository

import "errors"

var (
	ErrTeamNotFound     = errors.New("resource not found")
	ErrPrNotFound       = errors.New("resource not found")
	ErrReviewerNotFound = errors.New("resource not found")
	ErrAuthorNotFound   = errors.New("resource not found")
	ErrPrAlreadyExists  = errors.New("PR id already exists")
)
