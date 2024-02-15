package index

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/eid"
)

type IndexHandler struct {
	Eid *eid.EidFactory
}

func NewIndexHandler(l *slog.Logger, e *eid.EidFactory) IndexHandler {
	return IndexHandler{
		Eid: e,
	}
}

func (i IndexHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	component := index(i.Eid)
	component.Render(r.Context(), w)
}

func (i IndexHandler) Router(r chi.Router) {
	r.Get("/", i.indexHandler)
}
