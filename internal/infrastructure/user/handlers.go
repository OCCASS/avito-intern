package user

import (
	errorDto "github.com/OCCASS/avito-intern/internal/application/error"
	userDto "github.com/OCCASS/avito-intern/internal/application/user"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	"github.com/OCCASS/avito-intern/internal/domain/user"
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"
	"github.com/gofiber/fiber/v2"
)

type UserHandlers struct {
	services user.UserServices
}

func NewUserHandlers(s user.UserServices) *UserHandlers {
	return &UserHandlers{services: s}
}

func (h UserHandlers) SetIsActive(c *fiber.Ctx) error {
	payload := new(userDto.SetIsActiveDto)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("INVALID_BODY", "request body is invalid"))
	}

	user, err := h.services.SetIsActive(payload.UserId, payload.IsActive)
	if err != nil {
		return h.handleSetIsActiveError(c, err)
	}
	teamName, err := h.services.GetUserTeamName(payload.UserId)
	if err != nil {
		return h.handleSetIsActiveError(c, err)
	}
	return c.JSON(userDto.SetIsActiveResponse{
		User: userDto.SetIsActiveUserDetail{
			Id:       user.Id,
			Name:     user.Name,
			IsActive: user.IsActive,
			Team:     teamName,
		},
	})
}

func (h UserHandlers) handleSetIsActiveError(c *fiber.Ctx, err error) error {
	switch err {
	case userRepository.ErrUserNotFound, teamRepository.ErrTeamNotFound:
		return c.Status(fiber.StatusNotFound).JSON(errorDto.NewErrorResponse("NOT_FOUND", err.Error()))

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}

func (h UserHandlers) GetReview(c *fiber.Ctx) error {
	query := new(userDto.GetReviewQuery)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("QUERY_NOT_SET", "query params not set"))
	}

	prs, err := h.services.GetUserPullRequestsWhereReview(query.UserId)
	if err != nil {
		return h.handleGetReviewError(c, err)
	}
	return c.JSON(userDto.GetReviewResponse{UserId: query.UserId, PullRequests: prs})
}

func (h UserHandlers) handleGetReviewError(c *fiber.Ctx, err error) error {
	switch err {
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}
