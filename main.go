package main

import (
	"net/http"
	"log"
	"os"
	"context"
	"time"

	"github.com/akrylysov/algnhsa"
	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personal_website/index"
	"github.com/will-lol/personal_website/api"
) 

func url_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "url", r.URL)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	router := chi.NewRouter()
	router.Use(url_middleware)
	router.Route("/", index.Router) 
	router.Route("/api", api.Router)

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
