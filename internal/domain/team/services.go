package team

import (
	"github.com/OCCASS/avito-intern/internal/application/team"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"
	"github.com/OCCASS/avito-intern/internal/entity"
)

type TeamServices struct {
	teamRepository teamRepository.TeamRepository
	userRepository userRepository.UserRepository
}

func NewTeamServices(
	tr teamRepository.TeamRepository,
	ur userRepository.UserRepository,
) TeamServices {
	return TeamServices{
		teamRepository: tr,
		userRepository: ur,
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
	team, err := s.teamRepository.DeactivateMembers(dto.Name)
	if err != nil {
		return entity.Team{}, err
	}
	return team, nil
}
