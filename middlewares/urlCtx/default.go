package urlCtx

import (
	"context"
	"net/http"
)


func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "url", r.URL)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
