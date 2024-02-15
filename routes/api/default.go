package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/will-lol/personalWebsiteAwesome/db"
	"github.com/will-lol/personalWebsiteAwesome/routes/api/notifications"
	notificationsService "github.com/will-lol/personalWebsiteAwesome/services/notifications"
)

type ApiHandler struct {
	Logger *slog.Logger
	DB     *db.DB[notificationsService.Subscription]
}

func NewApiHandler(l *slog.Logger, d *db.DB[notificationsService.Subscription]) ApiHandler {
	return ApiHandler{
		Logger: l,
		DB:     d,
	}
}

func (a ApiHandler) Router(r chi.Router) {
	notifService, err := notificationsService.NewNotificationsService(a.Logger, a.DB)
	if err != nil {
		a.Logger.Error(err.Error())
	}
	notifHandler := notifications.NewNotificationsHandler(a.Logger, *notifService)
	r.Route("/notifications", notifHandler.Router)
}
