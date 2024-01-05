package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pdstuber/isit-a-cat/internal/api/handlers/healthcheck"
	"github.com/pdstuber/isit-a-cat/internal/api/handlers/imageupload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	*fiber.App
	listenPort string
	errChan    chan error
}

func NewRouter(config *Config) *Router {
	app := fiber.New()

	// TODO create all of the handlers, pass dependencies from constructor
	imageupload.NewHandler()
	app.Post("/images", metricsMiddleware.Decorate("postImageHandler", postImageHandler))
	app.Get("/predictions/{id}", metricsMiddleware.Decorate("getPredictionHandler", getPredictionHandler.ServeHTTP))
	app.Get("/images/{id}", metricsMiddleware.Decorate("getImageHandler", getImageHandler.ServeHTTP))
	app.Get("/ping", healthcheck.HandleHealth)
	app.Get("/metrics", promhttp.Handler())

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
