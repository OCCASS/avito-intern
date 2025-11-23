package repository

import "github.com/OCCASS/avito-intern/internal/entity"

type PullRequestRepository interface {
	Create(entity.PullRequest) (entity.PullRequest, error)
	Merge(entity.PullRequest) (entity.PullRequest, error)
	Reassign(pullRequestId, oldAuthor, newAuthor string) (entity.PullRequest, error)
}
