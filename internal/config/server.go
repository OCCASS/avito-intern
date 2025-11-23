package config

import (
	"github.com/gofiber/fiber/v2"
)

func NewServerConfig(config HTTPServerConfig) *fiber.Config {
	return &fiber.Config{
		ReadTimeout: config.ReadTimeout,
		IdleTimeout: config.IdleTimeout,
	}
}
