package server

import (
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	linkv1 "github.com/jr-dragon/dynamic_link/api/link/v1"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/phsym/console-slog"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/jr-dragon/dynamic_link/internal/data"
)

func NewHTTPServer(
	cfg data.Config,

	c *data.Clients,

	link *linkv1.Route,
) *fiber.App {
	app := fiber.New(fiber.Config{Prefork: cfg.HttpServer.Prefork})

	app.Use(slogfiber.New(logger(cfg)))
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: livenessProbe(c),
	}))

	link.RegisterHTTPRoutes(app)

	return app
}

func logger(cfg data.Config) *slog.Logger {
	var h slog.Handler
	if cfg.App.Logger.Pretty {
		h = console.NewHandler(os.Stderr, nil)
	} else {
		h = slog.NewTextHandler(os.Stderr, nil)
	}

	return slog.New(h)
}

func livenessProbe(clients *data.Clients) func(c *fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		return clients.RDB.Ping(c.Context()).Err() == nil
	}
}
