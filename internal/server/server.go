package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/OCCASS/avito-intern/internal/infrastructure/pullrequest"
	"github.com/OCCASS/avito-intern/internal/infrastructure/team"
	"github.com/OCCASS/avito-intern/internal/infrastructure/user"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app                 *fiber.App
	pullRequestHandlers *pullrequest.PullRequestHandlers
	teamHandlers        *team.TeamHandlers
	userHandlers        *user.UserHandlers
}

func NewServer(
	app *fiber.App,
	prh *pullrequest.PullRequestHandlers,
	th *team.TeamHandlers,
	uh *user.UserHandlers,
) *Server {
	return &Server{
		app:                 app,
		pullRequestHandlers: prh,
		teamHandlers:        th,
		userHandlers:        uh,
	}
}

func (s Server) SetupHandlers() {
	pullrequest := s.app.Group("/pullRequest")
	pullrequest.Post("/create", s.pullRequestHandlers.Create)
	pullrequest.Post("/merge", s.pullRequestHandlers.Merge)
	pullrequest.Post("/reassign", s.pullRequestHandlers.Reasign)

	team := s.app.Group("/team")
	team.Post("/add", s.teamHandlers.Add)
	team.Get("/get", s.teamHandlers.Get)

	user := s.app.Group("/users")
	user.Post("/setIsActive", s.userHandlers.SetIsActive)
	user.Get("/getReview", s.userHandlers.GetReview)
}

func (s *Server) MustStart(address string) {
	idleConnectionClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.app.Shutdown(); err != nil {
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnectionClosed)
	}()

	if err := s.app.Listen(address); err != nil {
		panic(fmt.Sprintf("Oops... Server is not running! Reason: %v", err))
	}

	<-idleConnectionClosed
}
