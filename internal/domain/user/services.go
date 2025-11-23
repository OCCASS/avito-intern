package user

import (
	pullrequestRepository "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"
	"github.com/OCCASS/avito-intern/internal/entity"
)

type UserServices struct {
	userRepository        userRepository.UserRepository
	teamRepository        teamRepository.TeamRepository
	pullrequestRepository pullrequestRepository.PullRequestRepository
}

func NewUserServices(
	ur userRepository.UserRepository,
	tr teamRepository.TeamRepository,
	pr pullrequestRepository.PullRequestRepository,
) UserServices {
	return UserServices{
		userRepository:        ur,
		teamRepository:        tr,
		pullrequestRepository: pr,
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

func (s UserServices) GetUserPullRequestsWhereReview(userId string) ([]entity.SmallPullRequest, error) {
	prs, err := s.pullrequestRepository.GetByReviewer(userId)
	if err != nil {
		return []entity.SmallPullRequest{}, err
	}
	return prs, nil
}
