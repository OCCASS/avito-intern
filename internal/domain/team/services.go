package team

import (
	pullrequestDto "github.com/OCCASS/avito-intern/internal/application/pullrequest"
	"github.com/OCCASS/avito-intern/internal/application/team"
	prServices "github.com/OCCASS/avito-intern/internal/domain/pullrequest"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	userServices "github.com/OCCASS/avito-intern/internal/domain/user"
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"

	"github.com/OCCASS/avito-intern/internal/entity"
)

type TeamServices struct {
	teamRepository teamRepository.TeamRepository
	userRepository userRepository.UserRepository

	pullrequestServices prServices.PullRequestServices
	userServices        userServices.UserServices
}

func NewTeamServices(
	tr teamRepository.TeamRepository,
	ur userRepository.UserRepository,
	prs prServices.PullRequestServices,
	us userServices.UserServices,
) TeamServices {
	return TeamServices{
		teamRepository:      tr,
		userRepository:      ur,
		pullrequestServices: prs,
		userServices:        us,
	}
}

func (s TeamServices) Add(dto team.CreateTeamDto) (entity.Team, error) {
	if err := s.userRepository.CreateMany(dto.Members); err != nil {
		return entity.Team{}, err
	}

	newTeam := entity.Team{
		Name:    dto.Name,
		Members: dto.Members,
	}
	team, err := s.teamRepository.Create(newTeam)
	if err != nil {
		return entity.Team{}, err
	}
	return team, nil
}

func (s TeamServices) Get(name string) (entity.Team, error) {
	team, err := s.teamRepository.Get(name)
	if err != nil {
		return entity.Team{}, err
	}
	return team, nil
}

func (s TeamServices) DeactivateMembers(dto team.DeactivateMembersDto) (entity.Team, error) {
	team, err := s.teamRepository.DeactivateMembers(dto.Name, dto.MembersIds)
	if err != nil {
		return entity.Team{}, err
	}

	for _, member := range team.Members {
		if member.IsActive {
			continue
		}

		prs, err := s.userServices.GetUserPullRequestsWhereReview(member.Id)
		if err != nil {
			continue
		}
		for _, pr := range prs {
			reassignDto := pullrequestDto.ReassignPullRequestDto{
				PullRequestId: pr.Id,
				OldReviewerId: member.Id,
				AllowRemove:   true,
			}
			_, _, err := s.pullrequestServices.Reassign(reassignDto)
			if err != nil {
				continue
			}
		}
	}

	team, err = s.teamRepository.Get(dto.Name)
	if err != nil {
		return entity.Team{}, err
	}
	return team, nil
}
