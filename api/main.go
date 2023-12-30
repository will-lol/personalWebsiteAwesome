package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/will-lol/personal_website/api/notifications"
)

func Router(r chi.Router) {
	r.Route("/notifications", notifications.Router)
}
