package link

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type (
	ApplinkHSet struct {
		Default     string `json:"default,omitempty" redis:"Default"`
		IOSLink     string `json:"ios,omitempty" redis:"iOS"`
		AndroidLink string `json:"android,omitempty" redis:"Android"`
		WindowsLink string `json:"windows,omitempty" redis:"Windows"`
		MacOSLink   string `json:"macos,omitempty" redis:"macOS"`
		LinuxLink   string `json:"linux,omitempty" redis:"Linux"`
	}

	CreateRequestApp struct {
		IOSLink     string `json:"ios,omitempty" redis:"iOS"`
		AndroidLink string `json:"android,omitempty" redis:"Android"`
		WindowsLink string `json:"windows,omitempty" redis:"Windows"`
		MacOSLink   string `json:"macos,omitempty" redis:"macOS"`
		LinuxLink   string `json:"linux,omitempty" redis:"Linux"`
	}

	CreateRequest struct {
		URL string            `json:"url" validate:"required,url"`
		App *CreateRequestApp `json:"app,omitempty"`
	}

	CreateResponse struct {
		URL string `json:"url"`
	}
)

func (link *Link) Create(ctx context.Context, req CreateRequest) (resp CreateResponse, err error) {
	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(req); err != nil {
		return
	}

	if req.IsAppLink() {
		return link.createApplink(ctx, req)
	}

	return link.createSimple(ctx, req)
}

func (link *Link) createSimple(ctx context.Context, req CreateRequest) (resp CreateResponse, err error) {
	p := link.rand.String(6, link.cfg.App.Key)

	resp.URL = link.cfg.App.RedirectorHost + "/s/" + p
	err = link.c.RDB.Set(ctx, p, req.URL, 0).Err()

	return resp, err
}

func (link *Link) createApplink(ctx context.Context, req CreateRequest) (resp CreateResponse, err error) {
	p := link.rand.String(6, link.cfg.App.Key)

	resp.URL = link.cfg.App.RedirectorHost + "/a/" + p
	err = link.c.RDB.HSet(ctx, p, req.ToApplinkHSet()).Err()

	return
}

func (req CreateRequest) IsSimple() bool {
	return req.App == nil
}

func (req CreateRequest) IsAppLink() bool {
	return req.App != nil
}

func (req CreateRequest) ToApplinkHSet() *ApplinkHSet {
	return &ApplinkHSet{
		Default:     req.URL,
		IOSLink:     req.App.IOSLink,
		AndroidLink: req.App.AndroidLink,
		WindowsLink: req.App.WindowsLink,
		MacOSLink:   req.App.MacOSLink,
		LinuxLink:   req.App.LinuxLink,
	}
}
