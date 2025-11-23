package repository

import "github.com/OCCASS/avito-intern/internal/entity"

type PullRequestRepository interface {
	Create(entity.PullRequest) (entity.PullRequest, error)
	Merge(string) (entity.PullRequest, error)
	Reassign(pullRequestId, oldReviewer, newReviewer string) (entity.PullRequest, error)
	Get(string) (entity.PullRequest, error)
}
