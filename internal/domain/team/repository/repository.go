package repository

import "github.com/OCCASS/avito-intern/internal/entity"

type TeamRepository interface {
	Create(entity.Team) (entity.Team, error)
	Get(string) (entity.Team, error)
	GetByUser(string) (entity.Team, error)
	DeactivateMembers(string, []string) (entity.Team, error)
}
