package pullrequest

import (
	errorDto "github.com/OCCASS/avito-intern/internal/application/error"
	pullrequestDto "github.com/OCCASS/avito-intern/internal/application/pullrequest"
	"github.com/OCCASS/avito-intern/internal/domain/pullrequest"
	pullrequestRepository "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	"github.com/gofiber/fiber/v2"
)

type PullRequestHandlers struct {
	services pullrequest.PullRequestServices
}

func NewPullRequestHandlers(s pullrequest.PullRequestServices) *PullRequestHandlers {
	return &PullRequestHandlers{services: s}
}

func (h PullRequestHandlers) Create(c *fiber.Ctx) error {
	payload := new(pullrequestDto.CreatePullRequestDto)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("INVALID_BODY", "request body is invalid"))
	}

	pr, err := h.services.Create(*payload)
	if err != nil {
		return h.handleCreateError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(pullrequestDto.CreatePullRequestResponse{Pr: pr})
}

func (h PullRequestHandlers) handleCreateError(c *fiber.Ctx, err error) error {
	switch err {
	case pullrequestRepository.ErrPrAlreadyExists:
		return c.Status(fiber.StatusConflict).JSON(errorDto.NewErrorResponse("PR_EXISTS", err.Error()))

	case teamRepository.ErrTeamNotFound,
		pullrequestRepository.ErrAuthorNotFound,
		pullrequestRepository.ErrTeamNotFound:
		return c.Status(fiber.StatusNotFound).JSON(errorDto.NewErrorResponse("NOT_FOUND", err.Error()))

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}

func (h PullRequestHandlers) Merge(c *fiber.Ctx) error {
	payload := new(pullrequestDto.MergePullRequestDto)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("INVALID_BODY", "request body is invalid"))
	}

	pr, err := h.services.Merge(*payload)
	if err != nil {
		return h.handleMergeError(c, err)
	}
	return c.JSON(pullrequestDto.MergePullRequestResponse{Pr: pr, MergedAt: *pr.MergedAt})
}

func (h PullRequestHandlers) handleMergeError(c *fiber.Ctx, err error) error {
	switch err {
	case pullrequestRepository.ErrPrNotFound:
		return c.Status(fiber.StatusNotFound).JSON(errorDto.NewErrorResponse("NOT_FOUND", err.Error()))

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}

func (h PullRequestHandlers) Reasign(c *fiber.Ctx) error {
	payload := new(pullrequestDto.ReassignPullRequestDto)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("INVALID_BODY", "request body is invalid"))
	}

	pr, replacedBy, err := h.services.Reassign(*payload)
	if err != nil {
		return h.handleReasignError(c, err)
	}
	return c.JSON(pullrequestDto.ReassignPullRequestResponse{Pr: pr, ReplacedBy: replacedBy})
}

func (h PullRequestHandlers) handleReasignError(c *fiber.Ctx, err error) error {
	switch err {
	case pullrequestRepository.ErrReviewerNotFound, pullrequestRepository.ErrPrNotFound:
		return c.Status(fiber.StatusNotFound).JSON(errorDto.NewErrorResponse("NOT_FOUND", err.Error()))

	case pullrequest.ErrPrMerged:
		return c.Status(fiber.StatusConflict).JSON(errorDto.NewErrorResponse("PR_MERGED", err.Error()))

	case pullrequest.ErrNoCandidatesToReassign:
		return c.Status(fiber.StatusConflict).JSON(errorDto.NewErrorResponse("NO_CANDIDATE", err.Error()))

	case pullrequest.ErrUserIsNotReviewer:
		return c.Status(fiber.StatusConflict).JSON(errorDto.NewErrorResponse("NOT_ASSIGNED", err.Error()))

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}
