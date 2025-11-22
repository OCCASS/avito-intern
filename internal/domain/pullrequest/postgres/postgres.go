package postgres

import (
	"github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/entity"
	"github.com/lib/pq"
)

type PullRequestPostgresRepository struct {
	db *database.Database
}

func NewPullRequestPostgresRepository(db *database.Database) *PullRequestPostgresRepository {
	return &PullRequestPostgresRepository{db}
}

func (r *PullRequestPostgresRepository) Create(pr entity.PullRequest) (entity.PullRequest, error) {
	newPr := &entity.PullRequest{}

	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return entity.PullRequest{}, err
	}

	queryPr := `INSERT INTO pullrequest(id, name, author_id, status) VALUES ($1, $2, $3, $4) RETURNING id, name, author_id, status`
	if err := tx.Get(&newPr, queryPr, pr.Id, pr.Name, pr.AuthorId, pr.Status); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}

	queryReviewers := `INSERT INTO pullrequest_reviewer (pullrequest_id, reviewer_id) SELECT $1 unset($2)`
	if _, err := tx.Exec(queryReviewers, pr.Id, pq.Array(pr.ReviewersIds)); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}
	newPr.ReviewersIds = pr.ReviewersIds

	if err := tx.Commit(); err != nil {
		return entity.PullRequest{}, err
	}

	return *newPr, nil
}

func (r *PullRequestPostgresRepository) SetStatus(id string, status entity.Status) (entity.PullRequest, error) {
	newPr := &entity.PullRequest{}

	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return entity.PullRequest{}, err
	}

	queryUpdateStatus := `UPDATE pullrequest SET status=$1 WHERE id=$2 RETURNING id, name, author_id, status`
	if err := tx.Get(&newPr, queryUpdateStatus, status, id); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}

	var reviewerIds []string
	queryReviewers := `SELECT reviewer_id FROM pullrequest_reviewer WHERE pullrequest_id=$1`
	if err := tx.Select(&reviewerIds, queryReviewers, id); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}

	newPr.ReviewersIds = reviewerIds

	if err := tx.Commit(); err != nil {
		return entity.PullRequest{}, err
	}

	return *newPr, nil
}

func (r *PullRequestPostgresRepository) Reassign(pullRequestId, oldAuthor, newAuthor string) (entity.PullRequest, error) {
	newPr := &entity.PullRequest{}

	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return entity.PullRequest{}, err
	}

	queryUpdateReviewer := `UPDATE pullrequest_reviewer SET reviewer_id=$1 WHERE pullrequest_id=$2 AND reviewer_id=$3`
	if _, err := tx.Exec(queryUpdateReviewer, newAuthor, pullRequestId, oldAuthor); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}

	queryPr := `SELECT
		p.id, p.name, p.author_id, p.status, array_agg(prr.reviewer_id) as reviewers
	FROM pullrequest p
	LEFT JOIN pullrequest_reviewers prr ON p.id = prr.pullrequest_id
	WHERE id=$1
	GROUP BY p.id`
	if err := tx.Get(&newPr, queryPr, pullRequestId); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.PullRequest{}, err
	}

	return *newPr, nil
}
