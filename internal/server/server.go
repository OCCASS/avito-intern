package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/OCCASS/avito-intern/internal/infrastructure/pullrequest"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app                 *fiber.App
	pullRequestHandlers *pullrequest.PullRequestHandlers
}

func NewServer(
	app *fiber.App,
	prh *pullrequest.PullRequestHandlers,
) *Server {
	return &Server{
		app:                 app,
		pullRequestHandlers: prh,
	}
}

func (s Server) SetupHandlers() {
	pullrequest := s.app.Group("/pullrequest")
	pullrequest.Post("/create", s.pullRequestHandlers.Create)
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
