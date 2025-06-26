package link

import (
	"context"

	"github.com/jr-dragon/dynamic_link/internal/biz/link/internal"
	"github.com/jr-dragon/dynamic_link/internal/data"
)

//go:generate moq -out mock_link.go . Contract

type Contract interface {
	CreateSimple(ctx context.Context, req CreateRequest) (resp CreateResponse, err error)
	ValidateSimple(code []byte) error
	RedirectSimple(ctx context.Context, code string) (url string, err error)
}

var _ Contract = &Link{}

type Link struct {
	cfg data.Config

	c *data.Clients

	rand internal.RandGenerator
}

func NewLink(cfg data.Config, c *data.Clients) Contract {
	return &Link{cfg: cfg, c: c, rand: &internal.Rand{}}
}
