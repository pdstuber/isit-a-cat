package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	*fiber.App
	listenPort string
	errChan    chan error
}

func NewRouter(config *Config) *Router {
	app := fiber.New()

	return &Router{
		App:        app,
		listenPort: config.ListenPort,
		errChan:    make(chan error),
	}
}

func (r *Router) Start(ctx context.Context) error {
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

func (r *Router) Stop(timeout time.Duration) error {
	return r.ShutdownWithTimeout(timeout)
}
