package pullrequest

import (
	"slices"

	"github.com/OCCASS/avito-intern/internal/application/pullrequest"
	prRepository "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	"github.com/OCCASS/avito-intern/internal/entity"

	"math/rand/v2"
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

	// TODO: refactor this part. move reviewers assign to another place
	filteredMembersIds := make([]string, 0, len(team.Members))
	for i := 0; i < len(team.Members); i++ {
		member := team.Members[i]
		if member.Id != dto.AuthorId && member.IsActive {
			filteredMembersIds = append(filteredMembersIds, member.Id)
		}
	}

	reviewers := make([]string, 0, 2)
	if len(filteredMembersIds) < 2 {
		reviewers = filteredMembersIds
	} else {
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

func (s PullRequestServices) Merge(dto pullrequest.MergePullRequestDto) (entity.PullRequest, error) {
	pr, err := s.pullRequestRepository.Merge(dto.Id)
	if err != nil {
		return entity.PullRequest{}, err
	}
	return pr, nil
}

func (s PullRequestServices) Reassign(dto pullrequest.ReassignPullRequestDto) (entity.PullRequest, *string, error) {
	pr, err := s.pullRequestRepository.Get(dto.PullRequestId)
	if err != nil {
		return entity.PullRequest{}, nil, err
	}

	if pr.Status == entity.StatusMerged {
		return entity.PullRequest{}, nil, ErrPrMerged
	} else if !slices.Contains(pr.ReviewersIds, dto.OldReviewerId) {
		return entity.PullRequest{}, nil, ErrUserIsNotReviewer
	}

	team, err := s.teamRepository.GetByUser(dto.OldReviewerId)
	if err != nil {
		return entity.PullRequest{}, nil, err
	}

	filteredMembersIds := make([]string, 0, len(team.Members))
	for i := 0; i < len(team.Members); i++ {
		member := team.Members[i]
		if member.Id != pr.AuthorId && !slices.Contains(pr.ReviewersIds, member.Id) && member.IsActive {
			filteredMembersIds = append(filteredMembersIds, member.Id)
		}
	}

	if len(filteredMembersIds) == 0 {
		return entity.PullRequest{}, nil, ErrNoCandidatesToReassign
	}

	newReviewerIndex := rand.IntN(len(filteredMembersIds))
	newReviewerId := filteredMembersIds[newReviewerIndex]
	pr, err = s.pullRequestRepository.Reassign(dto.PullRequestId, dto.OldReviewerId, newReviewerId)
	if err != nil {
		return entity.PullRequest{}, nil, err
	}
	return pr, &newReviewerId, nil

}
