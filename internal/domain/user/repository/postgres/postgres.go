package postgres

import (
	"database/sql"

	"github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/domain/user/repository"
	"github.com/OCCASS/avito-intern/internal/entity"
)

type UserPostgresRepository struct {
	db *database.Database
}

func NewUserPostgresRepository(db *database.Database) *UserPostgresRepository {
	return &UserPostgresRepository{db}
}

func (r UserPostgresRepository) UpdateIsActive(id string, newValue bool) (entity.User, error) {
	var user entity.User
	query := `UPDATE "user" SET is_active=$1 WHERE id=$2 RETURNING id, name, is_active`
	if err := r.db.Conn.Get(&user, query, newValue, id); err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, repository.ErrUserNotFound
		}
		return entity.User{}, nil
	}
	return user, nil
}

func (r UserPostgresRepository) GetReview(id string) ([]entity.PullRequest, error) {
	var prs []entity.PullRequest
	query := `SELECT
		pr.id,
		pr.name,
		pr.author_id,
		pr.status,
		pr.created_at,
		pr.merged_at
	FROM pullrequest_reviewer prr
	JOIN pullrequest pr
	ON pr.id = prr.pullrequest_id
	WHERE prr.reviewer_id=$1`
	if err := r.db.Conn.Select(&prs, query, id); err != nil {
		return []entity.PullRequest{}, err
	}
	return prs, nil
}
