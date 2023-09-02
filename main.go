package main

import (
	"fmt"
	"github.com/akrylysov/algnhsa"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"github.com/will-lol/personal_website/index"
	"time"
)

func main() {
	router := chi.NewRouter()
	router.Route("/", index.Router) 

	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {
		router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
		server := &http.Server{
			Addr:         "localhost:9000",
			Handler:      http.TimeoutHandler(router, 30*time.Second, "request timed out"),
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
		fmt.Printf("Listening on %v\n", server.Addr)
		server.ListenAndServe()
	} else {
		algnhsa.ListenAndServe(router, nil)
	}
}
