package notifications

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"

	"encoding/json"
	"io"
	"log"
)

type NotificationsHandler struct {
	Log                  *slog.Logger
	NotificationsService notifications.NotificationsService
}

func NewNotificationsHandler(l *slog.Logger, n notifications.NotificationsService) NotificationsHandler {
	return NotificationsHandler{
		Log:                  l,
		NotificationsService: n,
	}
}

func (n NotificationsHandler) Router(r chi.Router) {
	r.Post("/subscribe", n.subscribeHandler)
	r.Get("/notify", n.notifyHandler)
	r.Get("/publickey", n.publicKeyHandler)
}

func (n NotificationsHandler) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	subscriptionBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	var subscription notifications.Subscription
	json.Unmarshal(subscriptionBytes, &subscription)

	n.NotificationsService.Subscribe(subscription)
}

func (n NotificationsHandler) notifyHandler(w http.ResponseWriter, r *http.Request) {
	err := n.NotificationsService.Notify()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func (n NotificationsHandler) publicKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(n.NotificationsService.VapidPub))
}
