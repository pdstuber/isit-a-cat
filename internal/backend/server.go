package backend

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	*fiber.App
	listenPort string
	errChan    chan error
}

func NewServer(config *Config) *Server {
	app := fiber.New()

	return &Server{
		App:        app,
		listenPort: config.ListenPort,
		errChan:    make(chan error),
	}
}

func (r *Server) Start(ctx context.Context) error {
	go func() {
		if err := r.Listen(r.listenPort); err != nil {
			r.errChan <- err
		}
	}()

	select {
	case err := <-r.errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (r *Server) Stop(timeout time.Duration) error {
	return r.ShutdownWithTimeout(timeout)
}