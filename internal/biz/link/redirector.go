package link

import "context"

func (link *Link) RedirectSimple(ctx context.Context, code string) (string, error) {
	return link.c.RDB.Get(ctx, code).Result()
}
