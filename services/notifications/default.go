package notifications

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/will-lol/personalWebsiteAwesome/dependencies/db"
)

type Subscription = webpush.Subscription

type notificationsService interface {
	Notify() (err error)
	Subscribe(sub Subscription) (err error)
}

type NotificationsService struct {
	Log       *slog.Logger
	db        *db.DB[Subscription]
	vapidPub  *string
	vapidPriv *string
}

func NewNotificationsService(l *slog.Logger, d *db.DB[Subscription]) (*NotificationsService, error) {
	return &NotificationsService{
		Log:       l,
		db:        d,
	}, nil
}

func (n NotificationsService) getSecrets() (publicKey *string, privateKey *string, err error) {
	// Get required environment variables
	// The port of the parameters and secrets extension is needed because it may change
	port, err := strconv.Atoi(os.Getenv("PARAMETERS_SECRETS_EXTENSION_HTTP_PORT"))
	if err != nil {
		return nil, nil, err
	}

	// The AWS_SESSION_TOKEN is required to authenticate with the parameters and secrets manager
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	if len(sessionToken) == 0 {
		return nil, nil, errors.New("AWS_SESSION_TOKEN unset")
	}

	// The ARN of the secret is needed to request it from parameters and secrets manager
	arn := os.Getenv("SECRET_ARN")
	if len(arn) == 0 {
		return nil, nil, errors.New("SECRET_ARN unset")
	}

	// Make the request to parameters and secrets
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/secretsmanager/get?secretId=%s", port, arn), nil)
	req.Header.Add("X-AWS-Parameters-Secrets-Token", sessionToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New(resp.Status)
	}

	// Read the response bytes
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// The response is a JSON object that needs unmarshalling
	var rawUnmarshalling map[string]interface{}
	err = json.Unmarshal(bytes, &rawUnmarshalling)
	if err != nil {
		return nil, nil, err
	}
	secretString := rawUnmarshalling["SecretString"].(string) // Our secrets are hiding in this property of the object

	type vapidSecrets struct {
		PublicKey  string `json:"VAPID_PUB"`
		PrivateKey string `json:"VAPID_PRIV"`
	}
	var secrets vapidSecrets
	json.Unmarshal([]byte(secretString), &secrets)

	return &secrets.PublicKey, &secrets.PrivateKey, nil
}

func (n NotificationsService) GetPubKey() (key *string, err error) {
	if n.vapidPub != nil {
		return n.vapidPub, nil
	}
	n.vapidPub, n.vapidPriv, err = n.getSecrets()
	if err != nil {
		return nil, err 
	}
	return n.vapidPub, nil
}

func (n NotificationsService) GetPrivKey() (key *string, err error) {
	if n.vapidPriv != nil {
		return n.vapidPriv, nil
	}
	n.vapidPub, n.vapidPriv, err = n.getSecrets()
	if err != nil {
		return nil, err 
	}
	return n.vapidPriv, nil
}

func (n NotificationsService) Notify() (err error) {
	subscribers, err := n.db.GetObjects()
	if err != nil {
		n.Log.Error(err.Error())
		return err
	}

	for _, subscriber := range *subscribers {
		go n.notifySubscriber(subscriber)
	}
	n.Log.Info("Notified subscribers")
	return
}

func (n NotificationsService) notifySubscriber(sub Subscription) error {
	pub, err := n.GetPubKey()
	priv, err := n.GetPrivKey()
	if err != nil {
		return err
	}
	resp, err := webpush.SendNotification([]byte("Notification received"), &sub, &webpush.Options{
		VAPIDPublicKey:  *pub,
		VAPIDPrivateKey: *priv,
		TTL:             64,
		Subscriber:      "will.bradshaw50@gmail.com",
	})
	if err != nil {
		return err
	}
	if resp.StatusCode == 410 {
		n.db.DeleteObject(sub)
	}
	return nil
}

func (n NotificationsService) Subscribe(sub Subscription) error {
	_, err := url.ParseRequestURI(sub.Endpoint)
	if len(sub.Endpoint) < 1 || err != nil {
		n.Log.Error("Improper URL")
		return errors.New("Improper URL")
	}

	alreadyExists, err := n.db.DoesObjExist(sub)
	if err != nil {
		n.Log.Error("Couldn't check subscription")
		return errors.New("Couldn't check subscription")
	}
	if *alreadyExists {
		n.Log.Error("Already exists")
		return errors.New("Already exists")
	}

	err = n.db.SaveObject(sub)
	if err != nil {
		n.Log.Error(err.Error())
		return err
	}
	return nil
}
