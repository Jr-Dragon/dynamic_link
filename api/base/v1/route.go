package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jr-dragon/dynamic_link/api/internal/response"
	"github.com/jr-dragon/dynamic_link/internal/data"
)

type Route struct {
	c *data.Clients
}

func NewRoute(c *data.Clients) *Route {
	return &Route{c: c}
}

func (r *Route) RegisterHttpRoutes(app fiber.Router) {
	app.Get("/", r.root)
}

func (r *Route) root(c *fiber.Ctx) error {
	resp := response.Err(r.c.RDB.Ping(c.Context()).Err())

	c.Status(resp.Code)
	return c.JSON(resp)
}
