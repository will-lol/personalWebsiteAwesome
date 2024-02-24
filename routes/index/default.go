package index

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/services/icon"
)

type indexHandler struct {
	Log *slog.Logger
	IconService *icon.IconFinder
}

func NewIndexHandler(l *slog.Logger) (*indexHandler, error) {
	icons, err := icon.NewIconFinder("https://raw.githubusercontent.com/devicons/devicon/master/devicon.json")
	if err != nil {
		return nil, err
	}
	return &indexHandler{
		Log: l,
		IconService: &icons,
	}, nil
}

func (i indexHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	component := index(*i.IconService)
	component.Render(r.Context(), w)
}

func (i indexHandler) Router(r chi.Router) {
	r.Get("/", i.indexHandler)
}
