package stats

import (
	errorDto "github.com/OCCASS/avito-intern/internal/application/error"
	"github.com/OCCASS/avito-intern/internal/domain/stats"
	"github.com/gofiber/fiber/v2"
)

type StatsHandlers struct {
	services stats.StatsServices
}

func NewStatsHandlers(s stats.StatsServices) *StatsHandlers {
	return &StatsHandlers{services: s}
}

func (h StatsHandlers) Users(c *fiber.Ctx) error {
	stat, err := h.services.UsersStats()
	if err != nil {
		return h.handleUsersError(c, err)
	}
	return c.JSON(stat)
}

func (h StatsHandlers) handleUsersError(c *fiber.Ctx, err error) error {
	switch err {
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorDto.NewErrorResponse("SERVER_ERROR", "an internal error occurred"))
	}
}
