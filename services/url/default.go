package url

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
)

func GetURL(ctx context.Context) (*url.URL, error) {
	if req, ok := ctx.Value("url").(*url.URL); ok {
		return req, nil
	}
	return nil, errors.New("couldn't get url")
}

func GetURLHandled(ctx context.Context) (*url.URL) {
	if req, ok := ctx.Value("url").(*url.URL); ok {
		return req
	}
	slog.Default().Error("couldn't get url")
	return nil
}
