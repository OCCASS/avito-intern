package main

import (
	"flag"

	"github.com/OCCASS/avito-intern/internal/config"
	// "github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/server"
	"github.com/gofiber/fiber/v3"
)

func main() {
	cfgPath := flag.String("c", "", "Configuration file path.")
	flag.Parse()

	cfg := config.MustLoad(*cfgPath)
	serverCfg := config.NewServerConfig(cfg.Server)

	// db := database.MustConnect(cfg.Database)

	app := fiber.New(*serverCfg)
	httpServer := server.NewServer(app)
	httpServer.MustStart(cfg.Server.Address())
}
