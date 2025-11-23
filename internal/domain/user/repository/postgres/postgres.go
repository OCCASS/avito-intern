package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/domain/user/repository"
	"github.com/OCCASS/avito-intern/internal/entity"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

type UserPostgresRepository struct {
	db *database.Database
}

func NewUserPostgresRepository(db *database.Database) *UserPostgresRepository {
	return &UserPostgresRepository{db}
}

func (r UserPostgresRepository) CreateMany(users []entity.User) error {
	if len(users) == 0 {
		return nil
	}

	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO "user"(id, name, is_active) VALUES `
	vals := []interface{}{}
	placeholders := []string{}

	for i, u := range users {
		offset := i*3 + 1
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", offset, offset+1, offset+2))
		vals = append(vals, u.Id, u.Name, u.IsActive)
	}

	query += strings.Join(placeholders, ", ")
	if _, err = tx.Exec(query, vals...); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case pgerrcode.UniqueViolation:
				return repository.ErrUserAlreadyExists
			}
		}
		return err
	}

	return tx.Commit()
}

func (r UserPostgresRepository) UpdateIsActive(id string, newValue bool) (entity.User, error) {
	var user entity.User
	query := `UPDATE "user" SET is_active=$1 WHERE id=$2 RETURNING id, name, is_active`
	if err := r.db.Conn.Get(&user, query, newValue, id); err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, repository.ErrUserNotFound
		}
		return entity.User{}, err
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
