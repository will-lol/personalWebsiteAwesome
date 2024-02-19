package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/akrylysov/algnhsa"
	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/dependencies/db"
	"github.com/will-lol/personalWebsiteAwesome/dependencies/fs"
	"github.com/will-lol/personalWebsiteAwesome/dependencies/s3"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/eidCtx"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/envCtx"
	"github.com/will-lol/personalWebsiteAwesome/middlewares/urlCtx"
	"github.com/will-lol/personalWebsiteAwesome/routes"
	"github.com/will-lol/personalWebsiteAwesome/services/blog"
	"github.com/will-lol/personalWebsiteAwesome/services/env"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"
)

func main() {
	router := chi.NewRouter()
	router.Use(urlCtx.Middleware, envCtx.Middleware, eidCtx.Middleware)

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(l)

	var d db.DB[notifications.Subscription]
	if env.GetEnv(nil) != "dev" {
		var err error
		d, err = db.NewDB[notifications.Subscription]()
		if err != nil {
			l.Error("DB init failed")
		}
	}

	var f blog.Files
	if env.GetEnv(nil) == "dev" {
		f = fs.NewFs("blog")
	} else {
		var err error
		f, err = s3.NewS3()
		if err != nil {
			l.Error("S3 init failed")
		}
	}

	r, err := routes.NewRoutesHandler(l, d, &f)
	if err != nil {
		l.Error(err.Error())
	}
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
