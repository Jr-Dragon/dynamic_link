package link

import (
	"errors"
	"hash/crc32"
	"strings"
	"testing"
	"time"

	"github.com/gookit/goutil/testutil/assert"

	"github.com/jr-dragon/dynamic_link/internal/biz/link/internal"
	"github.com/jr-dragon/dynamic_link/internal/data"
	"github.com/jr-dragon/dynamic_link/internal/library/testutil"
)

const (
	code = "hSxIIn-v999y8"
	orig = "https://original.dynamic-link.test"
)

var cfg = data.Config{
	App: data.ConfigApp{
		RedirectorHost: "https://dynamic-link.test",
		Key:            crc32.IEEE,
	},
}

func TestLink_createSimple(t *testing.T) {
	testcases := []struct {
		name    string
		preHook func(c *data.Clients)
		req     CreateRequest
		wantErr error
	}{
		{
			name: "success",
			preHook: func(c *data.Clients) {
				c.RDBMock.
					ExpectSet(code, orig, time.Duration(0)).
					SetVal("OK")
			},
			req: CreateRequest{URL: orig},
		},
		{
			name:    "validation failed: missing url",
			req:     CreateRequest{},
			wantErr: errors.New("Key: 'CreateRequest.URL' Error:Field validation for 'URL' failed on the 'required' tag"),
		},
		{
			name:    "validation failed: invalid url",
			req:     CreateRequest{URL: "invalid-url"},
			wantErr: errors.New("Key: 'CreateRequest.URL' Error:Field validation for 'URL' failed on the 'url' tag"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := testutil.NewTestingClients()
			if err != nil {
				t.Fatal(err)
			}

			if tc.preHook != nil {
				tc.preHook(d)
			}

			link := NewLink(cfg, d).(*Link)
			link.rand = &internal.RandGeneratorMock{StringFunc: func(n int, k uint32) string { return code }}

			resp, err := link.createSimple(t.Context(), tc.req)
			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("CreateSimple() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr == nil {
				assert.True(t, strings.HasPrefix(resp.URL, cfg.App.RedirectorHost+"/s/"))
			}
		})
	}
}

func TestLink_createApplink(t *testing.T) {
	testcases := []struct {
		name    string
		preHook func(c *data.Clients)
		req     CreateRequest
		wantErr error
	}{
		{
			name: "success: with only default",
			preHook: func(c *data.Clients) {
				c.RDBMock.
					ExpectHSet(code, "Default", orig).
					SetVal(1)
			},
			req: CreateRequest{URL: orig, App: &CreateRequestApp{}},
		},
		{
			name: "success: with ios and android links",
			preHook: func(c *data.Clients) {
				c.RDBMock.
					ExpectHSet(code, "Default", orig, "Android", orig, "iOS", orig).
					SetVal(1)
			},
			req: CreateRequest{URL: orig, App: &CreateRequestApp{AndroidLink: orig, IOSLink: orig}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := testutil.NewTestingClients()
			if err != nil {
				t.Fatal(err)
			}

			if tc.preHook != nil {
				tc.preHook(d)
			}

			link := NewLink(cfg, d).(*Link)
			link.rand = &internal.RandGeneratorMock{StringFunc: func(n int, k uint32) string { return code }}

			resp, err := link.createApplink(t.Context(), tc.req)
			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("createApplink() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr == nil {
				assert.True(t, strings.HasPrefix(resp.URL, cfg.App.RedirectorHost+"/a/"))
			}
		})
	}
}
