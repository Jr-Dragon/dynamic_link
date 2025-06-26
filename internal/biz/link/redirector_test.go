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
