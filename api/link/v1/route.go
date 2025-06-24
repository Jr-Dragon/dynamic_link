package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jr-dragon/dynamic_link/api/internal/response"
	"github.com/jr-dragon/dynamic_link/internal/biz/link"
)

type Route struct {
	l *link.Link
}

func NewRoute(l *link.Link) *Route {
	return &Route{l: l}
}

func (r *Route) RegisterHTTPRoutes(app fiber.Router) {
	app.Post("link", r.create)
}

func (r *Route) create(c *fiber.Ctx) error {
	req := link.CreateRequest{}
	resp := &response.Response{}

	if err := c.BodyParser(&req); err != nil {
		resp = response.Err(err)
	}

	if d, err := r.l.Create(c.Context(), req); err != nil {
		resp = response.Err(err)
	} else {
		resp = response.Data(d)
	}

	c.Status(resp.Code)
	return c.JSON(resp)
}
