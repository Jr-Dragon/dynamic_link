package link

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jr-dragon/dynamic_link/internal/biz/link/internal"
)

type (
	CreateRequest struct {
		URL string `json:"url" validate:"required,url"`
	}

	CreateResponse struct {
		URL string `json:"url"`
	}
)

func (link *Link) CreateSimple(ctx context.Context, req CreateRequest) (resp CreateResponse, err error) {
	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(req); err != nil {
		return
	}

	p := internal.GenerateRandomString(6, link.cfg.App.Key)

	resp.URL = link.cfg.App.RedirectorHost + "/s/" + p
	err = link.c.RDB.Set(ctx, p, req.URL, 0).Err()

	return
}
