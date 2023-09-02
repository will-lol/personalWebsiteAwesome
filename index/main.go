package index

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func root(w http.ResponseWriter, r *http.Request) {
	component := index(time.Now().Local().Format(time.UnixDate))
	component.Render(r.Context(), w)
}

func Router(r chi.Router) {
	r.Get("/", root)
}
