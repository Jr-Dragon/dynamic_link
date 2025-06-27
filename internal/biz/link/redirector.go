package link

import (
	"context"

	"github.com/mileusna/useragent"
)

func (link *Link) RedirectSimple(ctx context.Context, code string) (string, error) {
	return link.c.RDB.Get(ctx, code).Result()
}

func (link *Link) RedirectApplink(ctx context.Context, uas, code string) (string, error) {
	ua := useragent.Parse(uas)

	res, err := link.c.RDB.HGetAll(ctx, code).Result()
	if err != nil {
		return "", err
	}

	if ua.IsAndroid() && res["Android"] != "" {
		return res["Android"], nil
	}
	if ua.IsIOS() && res["iOS"] != "" {
		return res["iOS"], nil
	}
	if ua.IsWindows() && res["Windows"] != "" {
		return res["Windows"], nil
	}
	if ua.IsMacOS() && res["macOS"] != "" {
		return res["macOS"], nil
	}
	if ua.IsLinux() && res["Linux"] != "" {
		return res["Linux"], nil
	}

	return res["Default"], nil
}
