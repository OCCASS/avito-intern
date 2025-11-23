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

	return c.JSON(pr)
}

/*
func (h PullRequestHandler) Merge(c *fiber.Ctx) error {

}

func (h PullRequestHandler) Reasign(c *fiber.Ctx) error {

}
*/
