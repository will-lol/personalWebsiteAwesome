package index

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	component := index("coolion")
	component.Render(r.Context(), w)
}

func Router(r chi.Router) {
	r.Get("/", root)
}
