package pullrequest

import (
	"crypto/sha256"
	"encoding/binary"
	"slices"
	"sort"

	"github.com/OCCASS/avito-intern/internal/application/pullrequest"
	prRepository "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	"github.com/OCCASS/avito-intern/internal/entity"
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

	reviewers := s.selectReviewers(team.Members, dto.AuthorId, dto.Id, 2)

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
	pr, err := s.pullRequestRepository.Get(dto.Id)
	if err != nil {
		return entity.PullRequest{}, err
	}

	if pr.Status == entity.StatusMerged {
		return pr, nil
	}

	return s.pullRequestRepository.Merge(dto.Id)
}

func (s PullRequestServices) Reassign(dto pullrequest.ReassignPullRequestDto) (entity.PullRequest, *string, error) {
	pr, err := s.pullRequestRepository.Get(dto.PullRequestId)
	if err != nil {
		return entity.PullRequest{}, nil, err
	}

	if err := s.validateReassignment(pr, dto.OldReviewerId); err != nil {
		return entity.PullRequest{}, nil, err
	}

	team, err := s.teamRepository.GetByUser(dto.OldReviewerId)
	if err != nil {
		return entity.PullRequest{}, nil, err
	}

	candidates := s.getReassignmentCandidates(team.Members, pr)
	newReviewerId, err := s.selectReviewer(dto.OldReviewerId, candidates)
	if err != nil {
		return entity.PullRequest{}, nil, err
	}

	pr, err = s.pullRequestRepository.Reassign(dto.PullRequestId, dto.OldReviewerId, newReviewerId)
	if err != nil {
		return entity.PullRequest{}, nil, err
	}

	return pr, &newReviewerId, nil
}

func (s PullRequestServices) selectReviewer(id string, canditates []string) (string, error) {
	if len(canditates) == 0 {
		return "", ErrNoCandidatesToReassign
	}

	sorted := s.deterministicSort(id, canditates)

	h := sha256.Sum256([]byte(id))
	num := binary.BigEndian.Uint64(h[:8])
	index := num % uint64(len(sorted))

	return sorted[index], nil
}

func (s PullRequestServices) selectReviewers(members []entity.User, authorId, oldReviewerId string, maxCount int) []string {
	validMembers := s.filterValidReviewers(members, authorId, nil)

	if len(validMembers) <= maxCount {
		return validMembers
	}

	validMembers = s.deterministicSort(oldReviewerId, validMembers)
	return validMembers[:maxCount]
}

func (s PullRequestServices) deterministicSort(id string, replacements []string) []string {
	sorted := make([]string, len(replacements))
	copy(sorted, replacements)

	sort.Slice(sorted, func(i, j int) bool {
		return s.hashPair(id, sorted[i]) < s.hashPair(id, sorted[j])
	})

	return sorted
}

func (s PullRequestServices) hashPair(id, value string) uint64 {
	h := sha256.Sum256([]byte(id + ":" + value))
	return binary.BigEndian.Uint64(h[:8])
}

func (s PullRequestServices) filterValidReviewers(members []entity.User, authorId string, excludeIds []string) []string {
	valid := make([]string, 0, len(members))

	for _, member := range members {
		if !member.IsActive || member.Id == authorId {
			continue
		}

		if excludeIds != nil && slices.Contains(excludeIds, member.Id) {
			continue
		}

		valid = append(valid, member.Id)
	}

	return valid
}

func (s PullRequestServices) getReassignmentCandidates(members []entity.User, pr entity.PullRequest) []string {
	return s.filterValidReviewers(members, pr.AuthorId, pr.ReviewersIds)
}

func (s PullRequestServices) validateReassignment(pr entity.PullRequest, oldReviewerId string) error {
	if pr.Status == entity.StatusMerged {
		return ErrPrMerged
	}

	if !slices.Contains(pr.ReviewersIds, oldReviewerId) {
		return ErrUserIsNotReviewer
	}

	return nil
}
