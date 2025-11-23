package pullrequest

import "github.com/OCCASS/avito-intern/internal/entity"

type CreatePullRequestDto struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
}

type MergePullRequestDto struct {
	Id string `json:"pull_request_id"`
}

type ReassignPullRequestDto struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerId string `json:"old_reviewer_id"`
}

type CreatePullRequestResponse struct {
	Pr entity.PullRequest `json:"pr"`
}

type MergePullRequestResponse struct {
	Pr       entity.PullRequest `json:"pr"`
	MergedAt string             `json:"mergedAt"`
}

type ReassignPullRequestResponse struct {
	Pr         entity.PullRequest `json:"pr"`
	ReplacedBy *string            `json:"replaced_by"`
}
