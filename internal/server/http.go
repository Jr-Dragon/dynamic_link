package server

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/phsym/console-slog"
	slogfiber "github.com/samber/slog-fiber"

	basev1 "github.com/jr-dragon/dynamic_link/api/base/v1"
	"github.com/jr-dragon/dynamic_link/internal/data"
)

func NewHTTPServer(
	cfg data.Config,

	base *basev1.Route,
) *fiber.App {
	app := fiber.New(fiber.Config{Prefork: cfg.HttpServer.Prefork})

	app.Use(slogfiber.New(logger(cfg)))
	//app.Use(otelfiber.Middleware())

	base.RegisterHttpRoutes(app)

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
