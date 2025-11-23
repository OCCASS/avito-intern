package main

import (
	"flag"

	"github.com/OCCASS/avito-intern/internal/config"
	"github.com/OCCASS/avito-intern/internal/database"
	"github.com/OCCASS/avito-intern/internal/server"
	"github.com/gofiber/fiber/v2"

	"github.com/OCCASS/avito-intern/internal/domain/pullrequest"
	prPostgres "github.com/OCCASS/avito-intern/internal/domain/pullrequest/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/domain/team"
	teamPostgres "github.com/OCCASS/avito-intern/internal/domain/team/repository/postgres"
	"github.com/OCCASS/avito-intern/internal/domain/user"
	userPostgres "github.com/OCCASS/avito-intern/internal/domain/user/repository/postgres"

	prHandlers "github.com/OCCASS/avito-intern/internal/infrastructure/pullrequest"
	tHandlers "github.com/OCCASS/avito-intern/internal/infrastructure/team"
	uHandlers "github.com/OCCASS/avito-intern/internal/infrastructure/user"
)

func main() {
	cfgPath := flag.String("c", "", "Configuration file path.")
	flag.Parse()

	cfg := config.MustLoad(*cfgPath)
	serverCfg := config.NewServerConfig(cfg.Server)

	db := database.MustConnect(cfg.Database)

	// Repositories
	pullRequestRepository := prPostgres.NewPullRequestPostgresRepository(db)
	teamRepository := teamPostgres.NewTeamPostgresRepository(db)
	userRepository := userPostgres.NewUserPostgresRepository(db)

	// Services
	pullRequestServices := pullrequest.NewPullRequestServices(pullRequestRepository, teamRepository)
	teamServices := team.NewTeamServices(teamRepository, userRepository)
	userServices := user.NewUserServices(userRepository, teamRepository, pullRequestRepository)

	// Handlers
	pullRequestHandlers := prHandlers.NewPullRequestHandlers(pullRequestServices)
	teamHandlers := tHandlers.NewTeamHandlers(teamServices)
	userHandlers := uHandlers.NewUserHandlers(userServices)

	app := fiber.New(*serverCfg)
	httpServer := server.NewServer(app, pullRequestHandlers, teamHandlers, userHandlers)
	httpServer.SetupHandlers()
	httpServer.MustStart(cfg.Server.Address())
}
