package routes

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/dependencies/db"
	"github.com/will-lol/personalWebsiteAwesome/routes/api"
	"github.com/will-lol/personalWebsiteAwesome/routes/blog"
	"github.com/will-lol/personalWebsiteAwesome/routes/index"
	blogService "github.com/will-lol/personalWebsiteAwesome/services/blog"
	"github.com/will-lol/personalWebsiteAwesome/services/env"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"
	"golang.org/x/sync/errgroup"
)

type RoutesHandler struct {
	Log         *slog.Logger
	DB          db.DB[notifications.Subscription]
	Files       *blogService.Files
	blogHandler RouteHandler
	indexHandler RouteHandler
	apiHandler  RouteHandler
}

type RouteHandler interface {
	Router(chi.Router)
}

func NewRoutesHandler(l *slog.Logger, d db.DB[notifications.Subscription], f *blogService.Files) (*RoutesHandler, error) {
	errs, _ := errgroup.WithContext(context.TODO())
	bHandlerCh := make(chan RouteHandler)
	iHandlerCh := make(chan RouteHandler)
	aHandlerCh := make(chan RouteHandler)

	errs.Go(func() error {
		handler, err := blog.NewBlogHandler(f, l) 
		if err != nil {
			l.Error(err.Error())
			return err
		}
		bHandlerCh <- handler
		return nil
	})
	errs.Go(func() error {
		handler, err := index.NewIndexHandler(l)
		if err != nil {
			l.Error(err.Error())
			return err
		}
		iHandlerCh <- handler
		return nil
	})
	go func() {
		handler := api.NewApiHandler(l, d)
		aHandlerCh <- handler
		return
	}()

	bHandler := <- bHandlerCh
	iHandler := <- iHandlerCh
	aHandler := <- aHandlerCh

	err := errs.Wait()
	if err != nil {
		return nil, err
	}

	return &RoutesHandler{
		Log:   l,
		DB:    d,
		Files: f,
		blogHandler: bHandler,
		indexHandler: iHandler,
		apiHandler: aHandler,
	}, nil
}

func (h RoutesHandler) Router(r chi.Router) {
	r.Route("/", h.indexHandler.Router)
	r.Route("/blog", h.blogHandler.Router)
	if env.GetEnv(nil) != "dev" {
		r.Route("/api", h.apiHandler.Router)
	}
}
