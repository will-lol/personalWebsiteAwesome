package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/db"
	"github.com/will-lol/personalWebsiteAwesome/eid"
	"github.com/will-lol/personalWebsiteAwesome/routes/api"
	"github.com/will-lol/personalWebsiteAwesome/routes/index"
	"github.com/will-lol/personalWebsiteAwesome/services/env"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"
)

type RoutesHandler struct {
	Log *slog.Logger
	DB  *db.DB[notifications.Subscription]
	Eid *eid.EidFactory
}

func NewRoutesHandler(l *slog.Logger, d *db.DB[notifications.Subscription], e *eid.EidFactory) RoutesHandler {
	return RoutesHandler{
		Log: l,
		DB:  d,
		Eid: e,
	}
}

func (h RoutesHandler) Router(r chi.Router) {
	r.Route("/", index.NewIndexHandler(h.Log, h.Eid).Router)
	if env.GetEnv(nil) != "dev" {
		r.Route("/api", api.NewApiHandler(h.Log, h.DB).Router)
	}
}
