package pullrequest

import (
	"github.com/OCCASS/avito-intern/internal/application/pullrequest"
	prRepository "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	"github.com/OCCASS/avito-intern/internal/entity"

	"math/rand"
	"time"
)

type PullRequestServices struct {
	pullRequestRepository prRepository.PullRequestRepository
	teamRepository        teamRepository.TeamRepository
}

func NewPullRequestServices(
	prR prRepository.PullRequestRepository,
	tR teamRepository.TeamRepository,
) PullRequestServices {
	return PullRequestServices{
		pullRequestRepository: prR,
		teamRepository:        tR,
	}
}

func (s PullRequestServices) Create(dto pullrequest.CreatePullRequestDto) (entity.PullRequest, error) {
	team, err := s.teamRepository.GetByUser(dto.AuthorId)
	if err != nil {
		return entity.PullRequest{}, err
	}

	filteredMembersIds := make([]string, 0, len(team.Members))
	for i := 0; i < len(team.Members); i++ {
		member := team.Members[i]
		if member.Id != dto.AuthorId {
			filteredMembersIds = append(filteredMembersIds, member.Id)
		}
	}

	reviewers := make([]string, 0, 2)
	if len(filteredMembersIds) < 2 {
		reviewers = filteredMembersIds
	} else {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(filteredMembersIds), func(i, j int) {
			filteredMembersIds[i], filteredMembersIds[j] = filteredMembersIds[j], filteredMembersIds[i]
		})

		reviewers = filteredMembersIds[:2]
	}

	pr := entity.PullRequest{
		Id:           dto.Id,
		Name:         dto.Name,
		AuthorId:     dto.AuthorId,
		ReviewersIds: reviewers,
		Status:       entity.StatusOpen,
	}

	return s.pullRequestRepository.Create(pr)
}
