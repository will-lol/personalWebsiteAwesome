package eid

import (
	"context"
	"errors"
	"log/slog"

	"github.com/will-lol/personalWebsiteAwesome/eid"
	"github.com/will-lol/personalWebsiteAwesome/lib/pointerify"
)

func GetEid(ctx context.Context) (*eid.EidFactory, error) {
	if req, ok := ctx.Value("eid").(*eid.EidFactory); ok {
		return req, nil
	}
	return nil, errors.New("couldn't get eid factory")
}

func GetNext(ctx context.Context) (*string, error) {
	factory, err := GetEid(ctx)
	if (err != nil) {
		return nil, err
	}
	return pointerify.Pointer(factory.GetNext()), nil
}

func GetNextHandled(ctx context.Context, l *slog.Logger) (string) {
	id, err := GetNext(ctx)
	if (err != nil) {
		l.Error(err.Error())
	}
	return *id
}
