package postgres

import (
	"github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/entity"
	"github.com/lib/pq"
)

type pullRequestRow struct {
	Id           string         `db:"id"`
	Name         string         `db:"name"`
	AuthorId     string         `db:"author_id"`
	Status       entity.Status  `db:"status"`
	CreatedAt    string         `db:"created_at"`
	MergedAt     *string        `db:"merged_at"`
	ReviewersIds pq.StringArray `db:"reviewers"`
}

type PullRequestPostgresRepository struct {
	db *database.Database
}

func NewPullRequestPostgresRepository(db *database.Database) *PullRequestPostgresRepository {
	return &PullRequestPostgresRepository{db}
}

func (r PullRequestPostgresRepository) Create(pr entity.PullRequest) (entity.PullRequest, error) {
	var newPr entity.PullRequest

	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return entity.PullRequest{}, err
	}
	defer tx.Rollback()

	queryPr := `INSERT INTO pullrequest(id, name, author_id, status, created_at)
	VALUES ($1, $2, $3, $4, now())
	RETURNING id, name, author_id, status, created_at`
	if err := tx.Get(&newPr, queryPr, pr.Id, pr.Name, pr.AuthorId, pr.Status); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}

	queryReviewers := `INSERT INTO pullrequest_reviewer(pullrequest_id, reviewer_id) SELECT $1, UNNEST($2::TEXT[])`
	if _, err := tx.Exec(queryReviewers, pr.Id, pq.Array(pr.ReviewersIds)); err != nil {
		tx.Rollback()
		return entity.PullRequest{}, err
	}
	newPr.ReviewersIds = pr.ReviewersIds

	if err := tx.Commit(); err != nil {
		return entity.PullRequest{}, err
	}

	return newPr, nil
}

func (r PullRequestPostgresRepository) Merge(id string) (entity.PullRequest, error) {
	var row pullRequestRow

	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return entity.PullRequest{}, err
	}
	defer tx.Rollback()

	queryMerge := `UPDATE pullrequest SET status=$1, merged_at=now() WHERE id=$2`
	if _, err := r.db.Conn.Exec(queryMerge, entity.StatusMerged, id); err != nil {
		return entity.PullRequest{}, err
	}

	queryPr := `SELECT
		p.id, p.name, p.author_id, p.status, array_agg(prr.reviewer_id) as reviewers, p.created_at, p.merged_at
	FROM pullrequest p
	LEFT JOIN pullrequest_reviewer prr ON p.id = prr.pullrequest_id
	WHERE id=$1
	GROUP BY p.id`
	if err := tx.Get(&row, queryPr, id); err != nil {
		return entity.PullRequest{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.PullRequest{}, err
	}

	pr := entity.PullRequest{
		Id:           row.Id,
		Name:         row.Name,
		AuthorId:     row.AuthorId,
		Status:       row.Status,
		ReviewersIds: []string(row.ReviewersIds),
		CreatedAt:    row.CreatedAt,
		MergedAt:     row.MergedAt,
	}

	return pr, nil
}

func (r PullRequestPostgresRepository) Reassign(pullRequestId, oldReviewer, newReviewer string) (entity.PullRequest, error) {
	var row pullRequestRow

	tx, err := r.db.Conn.Beginx()
	if err != nil {
		return entity.PullRequest{}, err
	}
	defer tx.Rollback()

	queryUpdateReviewer := `UPDATE pullrequest_reviewer SET reviewer_id=$1 WHERE pullrequest_id=$2 AND reviewer_id=$3`
	if _, err := tx.Exec(queryUpdateReviewer, newReviewer, pullRequestId, oldReviewer); err != nil {
		return entity.PullRequest{}, err
	}

	queryPr := `SELECT
		p.id, p.name, p.author_id, p.status, array_agg(prr.reviewer_id) as reviewers, p.created_at, p.merged_at
	FROM pullrequest p
	LEFT JOIN pullrequest_reviewer prr ON p.id = prr.pullrequest_id
	WHERE id=$1
	GROUP BY p.id`
	if err := tx.Get(&row, queryPr, pullRequestId); err != nil {
		return entity.PullRequest{}, err
	}

	if err := tx.Commit(); err != nil {
		return entity.PullRequest{}, err
	}

	pr := entity.PullRequest{
		Id:           row.Id,
		Name:         row.Name,
		AuthorId:     row.AuthorId,
		Status:       row.Status,
		ReviewersIds: []string(row.ReviewersIds),
		CreatedAt:    row.CreatedAt,
		MergedAt:     row.MergedAt,
	}

	return pr, nil
}

func (r PullRequestPostgresRepository) Get(id string) (entity.PullRequest, error) {
	var row pullRequestRow

	query := `SELECT
		p.id, p.name, p.author_id, p.status, array_agg(prr.reviewer_id) as reviewers, p.created_at, p.merged_at
	FROM pullrequest p
	LEFT JOIN pullrequest_reviewer prr ON p.id = prr.pullrequest_id
	WHERE id=$1
	GROUP BY p.id`
	if err := r.db.Conn.Get(&row, query, id); err != nil {
		return entity.PullRequest{}, err
	}

	pr := entity.PullRequest{
		Id:           row.Id,
		Name:         row.Name,
		AuthorId:     row.AuthorId,
		Status:       row.Status,
		ReviewersIds: []string(row.ReviewersIds),
		CreatedAt:    row.CreatedAt,
		MergedAt:     row.MergedAt,
	}
	return pr, nil
}
