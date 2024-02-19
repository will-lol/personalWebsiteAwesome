package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/dependencies/db"
	"github.com/will-lol/personalWebsiteAwesome/lib/pointerify"
	"github.com/will-lol/personalWebsiteAwesome/routes/api"
	"github.com/will-lol/personalWebsiteAwesome/routes/blog"
	"github.com/will-lol/personalWebsiteAwesome/routes/index"
	blogService "github.com/will-lol/personalWebsiteAwesome/services/blog"
	"github.com/will-lol/personalWebsiteAwesome/services/env"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"
)

type RoutesHandler struct {
	Log         *slog.Logger
	DB          db.DB[notifications.Subscription]
	Files       *blogService.Files
	blogHandler RouteHandler
	apiHandler  RouteHandler
}

type RouteHandler interface {
	Router(chi.Router)
}

func NewRoutesHandler(l *slog.Logger, d db.DB[notifications.Subscription], f *blogService.Files) (*RoutesHandler, error) {
	bHandler, err := blog.NewBlogHandler(f, l)
	if err != nil {
		return nil, err
	}

	return &RoutesHandler{
		Log:   l,
		DB:    d,
		Files: f,
		blogHandler: *bHandler,
		apiHandler: pointerify.DePointer(api.NewApiHandler(l, d)),
	}, nil
}

func (h RoutesHandler) Router(r chi.Router) {
	r.Route("/", index.Router)
	r.Route("/blog", h.blogHandler.Router)
	if env.GetEnv(nil) != "dev" {
		r.Route("/api", h.apiHandler.Router)
	}
}
