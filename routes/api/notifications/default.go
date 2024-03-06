package notifications

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personalWebsiteAwesome/services/notifications"

	"encoding/json"
	"io"
)

type notificationsHandler struct {
	Log                  *slog.Logger
	NotificationsService notifications.NotificationsService
}

func NewNotificationsHandler(l *slog.Logger, n notifications.NotificationsService) notificationsHandler {
	return notificationsHandler{
		Log:                  l,
		NotificationsService: n,
	}
}

func (n notificationsHandler) Router(r chi.Router) {
	r.Post("/subscribe", n.subscribeHandler)
	r.Get("/notify", n.notifyHandler)
	r.Get("/publickey", n.publicKeyHandler)
}

func (n notificationsHandler) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	subscriptionBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	var subscription notifications.Subscription
	json.Unmarshal(subscriptionBytes, &subscription)

	err = n.NotificationsService.Subscribe(subscription)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func (n notificationsHandler) notifyHandler(w http.ResponseWriter, r *http.Request) {
	err := n.NotificationsService.Notify()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func (n notificationsHandler) publicKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	pub, err := n.NotificationsService.GetPubKey()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.Write([]byte(*pub))
}
