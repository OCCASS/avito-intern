package pullrequest

import (
	pullrequestDto "github.com/OCCASS/avito-intern/internal/application/pullrequest"
	"github.com/OCCASS/avito-intern/internal/domain/pullrequest"
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
		return err
	}

	pr, err := h.services.Create(*payload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(pullrequestDto.CreatePullRequestResponse{Pr: pr})
}

func (h PullRequestHandlers) Merge(c *fiber.Ctx) error {
	payload := new(pullrequestDto.MergePullRequestDto)

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	pr, err := h.services.Merge(*payload)
	if err != nil {
		return err
	}
	return c.JSON(pullrequestDto.MergePullRequestResponse{Pr: pr, MergedAt: *pr.MergedAt})
}

func (h PullRequestHandlers) Reasign(c *fiber.Ctx) error {
	payload := new(pullrequestDto.ReassignPullRequestDto)

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	pr, replacedBy, err := h.services.Reassign(*payload)
	if err != nil {
		return err
	}
	return c.JSON(pullrequestDto.ReassignPullRequestResponse{Pr: pr, ReplacedBy: replacedBy})
}
