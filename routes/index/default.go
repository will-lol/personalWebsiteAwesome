package index

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	component := index()
	component.Render(r.Context(), w)
}

func Router(r chi.Router) {
	r.Get("/", indexHandler)
}
