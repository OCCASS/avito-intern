package team

import "github.com/OCCASS/avito-intern/internal/entity"

type CreateTeamDto struct {
	Name    string        `json:"team_name"`
	Members []entity.User `json:"members"`
}

type CreateTeamResponse struct {
	Team entity.Team `json:"team"`
}

type GetTeamQuery struct {
	Name string `query:"team_name"`
}
