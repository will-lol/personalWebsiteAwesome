package url

import (
	"net/url"
	"context"
	"errors"
)

func GetURL(ctx context.Context) (*url.URL, error) {
	if req, ok := ctx.Value("url").(*url.URL); ok {
		return req, nil
	}
	return nil, errors.New("couldn't get url")
}
