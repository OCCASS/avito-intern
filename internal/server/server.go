package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app *fiber.App
}

func NewServer(app *fiber.App) *Server {
	return &Server{
		app,
	}
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
