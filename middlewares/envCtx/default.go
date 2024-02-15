package envCtx

import (
	"context"
	"net/http"
	"os"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "env", os.Getenv("ENVIRONMENT"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
} 
