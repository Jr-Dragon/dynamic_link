package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/jr-dragon/dynamic_link/api/internal/response"
	"github.com/jr-dragon/dynamic_link/internal/biz/link"
)

func TestRoute_createSimple(t *testing.T) {
	testcases := []struct {
		name       string
		link       *link.ContractMock
		req        link.CreateRequest
		wantResp   *response.Response
		wantStatus int
	}{
		{
			name: "success",
			link: &link.ContractMock{CreateSimpleFunc: func(ctx context.Context, req link.CreateRequest) (link.CreateResponse, error) {
				return link.CreateResponse{URL: "https://dynamic-link.test/s/hSxIIn-v999y8"}, nil
			}},
			req:        link.CreateRequest{URL: "https://original.dynamic-link.test/"},
			wantResp:   response.Data(link.CreateResponse{URL: "https://dynamic-link.test/s/hSxIIn-v999y8"}),
			wantStatus: http.StatusOK,
		},
		{
			name: "validate failed",
			link: &link.ContractMock{CreateSimpleFunc: func(ctx context.Context, req link.CreateRequest) (link.CreateResponse, error) {
				return link.CreateResponse{URL: ""}, validator.ValidationErrors{}
			}},
			req:        link.CreateRequest{URL: "invalid_url_string"},
			wantResp:   response.Err(validator.ValidationErrors{}),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "database error",
			link: &link.ContractMock{CreateSimpleFunc: func(ctx context.Context, req link.CreateRequest) (link.CreateResponse, error) {
				return link.CreateResponse{URL: ""}, redis.ErrClosed
			}},
			req:        link.CreateRequest{URL: "https://original.dynamic-link.test/"},
			wantResp:   response.Err(redis.ErrClosed),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			app := setup(tc.link)
			req, _ := http.NewRequest("POST", "/link", bytes.NewReader(serialize(tc.req)))

			res, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("TestRoute_create() error = %v", err)
			}

			assert.Equal(t, tc.wantStatus, res.StatusCode)
			assert.JSONEq(t, string(serialize(tc.wantResp)), string(serialize(res.Body)))
		})
	}
}

func TestRoute_redirectSimple(t *testing.T) {
	testcases := []struct {
		name       string
		link       *link.ContractMock
		code       string
		wantResp   *response.Response
		wantStatus int
	}{
		{
			name: "success",
			link: &link.ContractMock{
				ValidateSimpleFunc: func(code []byte) error { return nil },
				RedirectSimpleFunc: func(ctx context.Context, code string) (string, error) {
					return "https://original.dynamic-link.test", nil
				},
			},
			code:       "hSxIIn-v999y8",
			wantStatus: http.StatusFound,
		},
		{
			name: "validate failed",
			link: &link.ContractMock{
				ValidateSimpleFunc: func(code []byte) error { return errors.New("") },
			},
			code:       "invalid_url_string",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			app := setup(tc.link)
			req, _ := http.NewRequest("GET", "/s/"+tc.code, bytes.NewReader(serialize(tc.code)))

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("TestRoute_redirect() error = %v", err)
			}

			assert.Equal(t, tc.wantStatus, resp.StatusCode)
		})
	}
}

func setup(l link.Contract) *fiber.App {
	app := fiber.New()
	NewRoute(l).RegisterHTTPRoutes(app)

	return app
}

func serialize(payload any) (r []byte) {
	switch payload.(type) {
	case io.Reader:
		r, _ = io.ReadAll(payload.(io.Reader))
	default:
		r, _ = json.Marshal(payload)
	}
	return r
}
