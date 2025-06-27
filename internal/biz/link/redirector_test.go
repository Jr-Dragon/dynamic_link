package link

import (
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/jr-dragon/dynamic_link/internal/data"
	"github.com/jr-dragon/dynamic_link/internal/library/testutil"
)

func TestLink_RedirectSimple(t *testing.T) {
	testcases := []struct {
		name    string
		preHook func(c *data.Clients)
		code    string
		want    string
		wantErr error
	}{
		{
			name:    "success",
			preHook: func(c *data.Clients) { c.RDBMock.ExpectGet(code).SetVal(orig) },
			code:    code,
			want:    orig,
		},
		{
			name:    "error: expired code",
			preHook: func(c *data.Clients) { c.RDBMock.ExpectGet(code).RedisNil() },
			code:    code,
			wantErr: redis.Nil,
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

			l := NewLink(cfg, d)
			got, err := l.RedirectSimple(t.Context(), tc.code)
			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("RedirectSimple() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err == nil && got != tc.want {
				t.Errorf("RedirectSimple() got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLink_RedirectApplink(t *testing.T) {
	testcases := []struct {
		name    string
		preHook func(c *data.Clients)
		uas     string
		code    string
		want    string
		wantErr error
	}{
		{
			name: "success",
			preHook: func(c *data.Clients) {
				c.RDBMock.ExpectHGetAll(code).SetVal(map[string]string{"Default": orig})
			},
			uas:  "",
			code: code,
			want: orig,
		},
		{
			name: "success: ios link",
			preHook: func(c *data.Clients) {
				c.RDBMock.ExpectHGetAll(code).SetVal(map[string]string{"Default": orig, "iOS": "https://apple.com/"})
			},
			uas:  "Mozilla/5.0 (iPad; CPU OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H321 Safari/600.1.4",
			code: code,
			want: "https://apple.com/",
		},
		{
			name: "success: fallback to default",
			preHook: func(c *data.Clients) {
				c.RDBMock.ExpectHGetAll(code).SetVal(map[string]string{"Default": orig, "Android": "https://android.com/"})
			},
			uas:  "Mozilla/5.0 (iPad; CPU OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H321 Safari/600.1.4",
			code: code,
			want: orig,
		},
		{
			name:    "error: expired code",
			preHook: func(c *data.Clients) { c.RDBMock.ExpectHGetAll(code).RedisNil() },
			code:    code,
			wantErr: redis.Nil,
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

			l := NewLink(cfg, d)
			got, err := l.RedirectApplink(t.Context(), tc.uas, tc.code)
			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("RedirectSimple() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err == nil && got != tc.want {
				t.Errorf("RedirectSimple() got = %v, want %v", got, tc.want)
			}
		})
	}
}
