package team

import (
	errorDto "github.com/OCCASS/avito-intern/internal/application/error"
	teamDto "github.com/OCCASS/avito-intern/internal/application/team"
	"github.com/OCCASS/avito-intern/internal/domain/team"
	teamRepository "github.com/OCCASS/avito-intern/internal/domain/team/repository"
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"
	"github.com/gofiber/fiber/v2"
)

type TeamHandlers struct {
	services team.TeamServices
}

func NewTeamHandlers(s team.TeamServices) *TeamHandlers {
	return &TeamHandlers{services: s}
}

func (h TeamHandlers) Add(c *fiber.Ctx) error {
	payload := new(teamDto.CreateTeamDto)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("INVALID_BODY", "request body is invalid"))
	}

	team, err := h.services.Add(*payload)
	if err != nil {
		return h.handleAddError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(teamDto.CreateTeamResponse{Team: team})
}

func (h TeamHandlers) handleAddError(c *fiber.Ctx, err error) error {
	switch err {
	case teamRepository.ErrTeamAlreadyExists:
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("TEAM_EXISTS", err.Error()))

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}

func (h TeamHandlers) Get(c *fiber.Ctx) error {
	query := new(teamDto.GetTeamQuery)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorDto.NewErrorResponse("QUERY_NOT_SET", "query params not set"))
	}

	team, err := h.services.Get(query.Name)
	if err != nil {
		return h.handleGetError(c, err)
	}

	return c.JSON(team)
}

func (h TeamHandlers) handleGetError(c *fiber.Ctx, err error) error {
	switch err {
	case teamRepository.ErrTeamNotFound:
		return c.Status(fiber.StatusNotFound).JSON(errorDto.NewErrorResponse("NOT_FOUND", err.Error()))
	case userRepository.ErrUserAlreadyExists:
		return c.Status(fiber.StatusConflict).JSON(errorDto.NewErrorResponse("USER_EXISTS", err.Error()))

	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}
