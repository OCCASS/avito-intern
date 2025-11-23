package repository

import "github.com/OCCASS/avito-intern/internal/entity"

type UserRepository interface {
	CreateMany([]entity.User) error
	UpdateIsActive(string, bool) (entity.User, error)
	GetReview(string) ([]entity.PullRequest, error)
}
