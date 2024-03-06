package url

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
)

const hostname = "w.bradshaw.page"

func GetURL(ctx context.Context) (*url.URL, error) {
	if req, ok := ctx.Value("url").(*url.URL); ok {
		req.Host = hostname
		req.Scheme = "https"
		return req, nil
	}
	return nil, errors.New("couldn't get url")
}

func GetURLHandled(ctx context.Context) *url.URL {
	val, err := GetURL(ctx)
	if err != nil {
		slog.Default().Error(err.Error())
	}
	return val
}
