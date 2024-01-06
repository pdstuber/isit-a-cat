package api

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pdstuber/isit-a-cat/internal/api/handlers/getprediction"
	"github.com/pdstuber/isit-a-cat/internal/api/handlers/imageretrieval"
	"github.com/pdstuber/isit-a-cat/internal/api/handlers/postimage"
	"github.com/pdstuber/isit-a-cat/internal/dep"
)

type Router struct {
	fiberApp   *fiber.App
	listenPort string
	errChan    chan error
}

type routerDependencies interface {
	dep.CanForwardDependencies
}

func NewRouter(deps routerDependencies, listenPort string) *Router {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(healthcheck.New())
	app.Use(cors.New())
	app.Use(recover.New())

	app.Use("/predictions", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	postImageHandler := postimage.NewHandler(deps.Forward())
	getPredictionHandler := getprediction.NewHandler(deps.Forward())
	getImageHandler := imageretrieval.NewHandler(deps.Forward())

	// TODO move bot to webhook and include here
	app.Post("/images", postImageHandler.Handle)
	app.Get("/predictions/:id", getPredictionHandler.Handle)
	app.Get("/images/:id", getImageHandler.Handle)

	return &Router{
		fiberApp:   app,
		listenPort: listenPort,
		errChan:    make(chan error),
	}
}

func (r *Router) Start(ctx context.Context) error {
	log.Println("starting router")
	go func() {
		if err := r.fiberApp.Listen(r.listenPort); err != nil {
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
	return r.fiberApp.ShutdownWithTimeout(timeout)
}
