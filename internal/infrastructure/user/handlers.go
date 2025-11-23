package user

import (
	"github.com/OCCASS/avito-intern/internal/domain/user"
	"github.com/gofiber/fiber/v2"
)

type UserHandlers struct {
	services user.UserServices
}

func NewUserHandlers(s user.UserServices) *UserHandlers {
	return &UserHandlers{services: s}
}

func (h UserHandlers) SetIsActive(c *fiber.Ctx) error { return nil }

func (h UserHandlers) GetReview(c *fiber.Ctx) error { return nil }
