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
	"github.com/will-lol/personalWebsiteAwesome/eid"
	"github.com/will-lol/personalWebsiteAwesome/routes"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/urlCtx"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/envCtx"
)

func main() {
	router := chi.NewRouter()
	router.Use(urlCtx.Middleware, envCtx.Middleware)

	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	d, err := db.NewDB[notifications.Subscription]()
	if (err != nil) {
		l.Warn("DB init failed")
	}
	e := eid.NewEidFactory()

	r := routes.NewRoutesHandler(l, d, &e)
	r.Router(router)

	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {
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
