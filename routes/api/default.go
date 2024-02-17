package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/will-lol/personalWebsiteAwesome/dependencies/db"
	"github.com/will-lol/personalWebsiteAwesome/routes/api/notifications"
	notificationsService "github.com/will-lol/personalWebsiteAwesome/services/notifications"
)

type ApiHandler struct {
	Log          *slog.Logger
	DB           *db.DB[notificationsService.Subscription]
	notificationsService *notificationsService.NotificationsService
}

func NewApiHandler(l *slog.Logger, d *db.DB[notificationsService.Subscription]) *ApiHandler {
	notifService, err := notificationsService.NewNotificationsService(l, d)
	if err != nil {
		l.Error(err.Error())
	}
	return &ApiHandler{
		Log: l,
		DB:  d,
		notificationsService: notifService,
	}
}

func (a ApiHandler) Router(r chi.Router) {
	notifHandler := notifications.NewNotificationsHandler(a.Log, a.notificationsService)
	r.Route("/notifications", notifHandler.Router)
}
