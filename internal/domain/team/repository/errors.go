package repository

import "errors"

var (
	ErrTeamAlreadyExists = errors.New("team_name already exists")
	ErrTeamNotFound      = errors.New("resource not found")
)
