package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/akrylysov/algnhsa"
	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/db"
	"github.com/will-lol/personalWebsiteAwesome/routes"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/urlCtx"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/envCtx"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/eidCtx"
	"github.com/will-lol/personalWebsiteAwesome/services/env"
)

func main() {
	router := chi.NewRouter()
	router.Use(urlCtx.Middleware, envCtx.Middleware, eidCtx.Middleware)

	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(l)
	d, err := db.NewDB[notifications.Subscription]()
	if (err != nil) {
		l.Warn("DB init failed")
	}

	r := routes.NewRoutesHandler(l, d)
	r.Router(router)

	if env.GetEnv(nil) == "dev" {
		router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
		server := &http.Server{
			Addr:         "localhost:9000",
			Handler:      http.TimeoutHandler(router, 30*time.Second, "request timed out"),
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		log.Printf("Listening on %v\n", server.Addr)
		server.ListenAndServe()
	} else {
		algnhsa.ListenAndServe(router, nil)
	}
}
