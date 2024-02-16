package eidCtx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/will-lol/personalWebsiteAwesome/eid"
)

func Middleware(next http.Handler) http.Handler {
	eid := eid.NewEidFactory()
	slog.Default().Info(eid.GetNext())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "eid", &eid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
