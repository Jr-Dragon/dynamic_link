package link

import (
	"github.com/jr-dragon/dynamic_link/internal/data"
)

type Link struct {
	cfg data.Config

	c *data.Clients
}

func NewLink(cfg data.Config, c *data.Clients) *Link {
	return &Link{cfg: cfg, c: c}
}
