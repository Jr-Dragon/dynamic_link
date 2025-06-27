package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/jr-dragon/dynamic_link/api/internal/response"
	"github.com/jr-dragon/dynamic_link/internal/biz/link"
)

type Route struct {
	l link.Contract
}

func NewRoute(l link.Contract) *Route {
	return &Route{l: l}
}

func (r *Route) RegisterHTTPRoutes(app fiber.Router) {
	app.Post("link", r.create)

	app.Get("s/:code", r.validateSimple, r.redirectSimple)
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

func (r *Route) validateSimple(c *fiber.Ctx) error {
	if err := r.l.ValidateSimple([]byte(c.Params("code"))); err != nil {
		resp := response.Err(response.InvalidCode(err))

		c.Status(resp.Code)
		return c.JSON(resp)
	}

	return c.Next()
}

func (r *Route) redirectSimple(c *fiber.Ctx) error {
	if url, err := r.l.RedirectSimple(c.Context(), c.Params("code")); err != nil {
		var resp *response.Response
		if errors.Is(err, redis.Nil) {
			resp = response.Err(response.ExpiredCode(err))
		} else {
			resp = response.Err(err)
		}

		c.Status(resp.Code)
		return c.JSON(resp)
	} else {
		return c.Redirect(url, http.StatusFound)
	}
}
