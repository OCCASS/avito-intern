package main

import (
	"flag"

	"github.com/OCCASS/avito-intern/internal/config"
	"github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/server"
	"github.com/gofiber/fiber/v2"

	"github.com/OCCASS/avito-intern/internal/domain/pullrequest"
	prPosgres "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository/postgres"
	teamPosgres "github.com/OCCASS/avito-intern/internal/domain/team/repository/postgres"

	prHandlers "github.com/OCCASS/avito-intern/internal/infrastructure/pullrequest"
)

func main() {
	cfgPath := flag.String("c", "", "Configuration file path.")
	flag.Parse()

	cfg := config.MustLoad(*cfgPath)
	serverCfg := config.NewServerConfig(cfg.Server)

	db := database.MustConnect(cfg.Database)

	pullRequestRepository := prPosgres.NewPullRequestPostgresRepository(db)
	teamRepository := teamPosgres.NewTeamPostgresRepository(db)

	pullRequestService := pullrequest.NewPullRequestServices(pullRequestRepository, teamRepository)

	pullRequestHandlers := prHandlers.NewPullRequestHandlers(pullRequestService)

	app := fiber.New(*serverCfg)
	httpServer := server.NewServer(app, pullRequestHandlers)
	httpServer.SetupHandlers()
	httpServer.MustStart(cfg.Server.Address())
}
