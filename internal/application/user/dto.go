package user

import "github.com/OCCASS/avito-intern/internal/entity"

type SetIsActiveDto struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveUserDetail struct {
	Id       string `json:"user_id"`
	Name     string `json:"username"`
	Team     string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	User SetIsActiveUserDetail `json:"user"`
}

type GetReviewQuery struct {
	UserId string `query:"user_id"`
}

type GetReviewResponse struct {
	UserId       string                    `json:"user_id"`
	PullRequests []entity.SmallPullRequest `json:"pull_requests"`
}
