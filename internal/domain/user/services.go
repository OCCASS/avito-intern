package user

import (
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"
	"github.com/OCCASS/avito-intern/internal/entity"
)

type UserServices struct {
	userRepository userRepository.UserRepository
	teamRepository teamRepository.TeamRepository
}

func NewUserServices(
	ur userRepository.UserRepository,
	tr teamRepository.TeamRepository,
) UserServices {
	return UserServices{
		userRepository: ur,
		teamRepository: tr,
	}
}

func (s UserServices) SetIsActive(userId string, isActive bool) (entity.User, error) {
	user, err := s.userRepository.UpdateIsActive(userId, isActive)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (s UserServices) GetUserTeamName(userId string) (string, error) {
	team, err := s.teamRepository.GetByUser(userId)
	if err != nil {
		return "", err
	}
	return team.Name, nil
}
