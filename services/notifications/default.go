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
	"sync"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/will-lol/personalWebsiteAwesome/dependencies/db"
)

type Subscription = webpush.Subscription

type NotificationsService interface {
	Notify() (err error)
	Subscribe(sub Subscription) (err error)
	GetPubKey() (*string, error)
	GetPrivKey() (*string, error)
}

type notificationsService struct {
	Log       *slog.Logger
	db        db.DB[Subscription]
	vapidPub  *string
	vapidPriv *string
}

func NewNotificationsService(l *slog.Logger, d db.DB[Subscription]) (*notificationsService, error) {
	return &notificationsService{
		Log: l,
		db:  d,
	}, nil
}

func (n notificationsService) getSecrets() (publicKey *string, privateKey *string, err error) {
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

func (n notificationsService) GetPubKey() (key *string, err error) {
	if n.vapidPub != nil {
		return n.vapidPub, nil
	}
	n.vapidPub, n.vapidPriv, err = n.getSecrets()
	if err != nil {
		return nil, err
	}
	return n.vapidPub, nil
}

func (n notificationsService) GetPrivKey() (key *string, err error) {
	if n.vapidPriv != nil {
		return n.vapidPriv, nil
	}
	n.vapidPub, n.vapidPriv, err = n.getSecrets()
	if err != nil {
		return nil, err
	}
	return n.vapidPriv, nil
}

func (n notificationsService) Notify() (err error) {
	subscribers, err := n.db.GetObjects()
	if err != nil {
		n.Log.Error(err.Error())
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(*subscribers))
	n.Log.Debug("notifying subscribers: ", "subscribers", *subscribers)
	for _, subscriber := range *subscribers {
		wg.Add(1)
		go func(sub Subscription) {
			defer wg.Done()
			err := n.notifySubscriber(sub)
			if err != nil {
				n.Log.Debug("there was an err: " + err.Error())
				errCh <- err
			}
		}(subscriber)
	}

	go func() {
		n.Log.Debug("waiting")
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		return err
	}

	wg.Wait()
	n.Log.Info("Notified subscribers")
	return
}

func (n notificationsService) notifySubscriber(sub Subscription) error {
	n.Log.Debug("notify subscriber")
	pub, err := n.GetPubKey()
	priv, err := n.GetPrivKey()
	if err != nil {
		return err
	}
	n.Log.Debug("got keys")
	n.Log.Debug("sending using webpush")
	resp, err := webpush.SendNotification([]byte("Notification received"), &sub, &webpush.Options{
		VAPIDPublicKey:  *pub,
		VAPIDPrivateKey: *priv,
		TTL:             64,
		Subscriber:      "will.bradshaw50@gmail.com",
	})
	if err != nil {
		return err
	}
	n.Log.Debug("sent")
	if resp.StatusCode == 410 {
		n.Log.Debug("deleting")
		n.db.DeleteObject(sub)
	}
	n.Log.Debug("returning after success...")
	return nil
}

func (n notificationsService) Subscribe(sub Subscription) error {
	n.Log.Debug("subscription", "subscription", sub)
	n.Log.Debug("checking url validity")
	_, err := url.ParseRequestURI(sub.Endpoint)
	if len(sub.Endpoint) < 1 || err != nil {
		n.Log.Error("Improper URL")
		return errors.New("Improper URL")
	}
	n.Log.Debug("url valid")

	n.Log.Debug("checking if obj exists")
	alreadyExists, err := n.db.DoesObjExist(sub)
	if err != nil {
		n.Log.Error("Couldn't check subscription")
		return errors.New("Couldn't check subscription")
	}
	n.Log.Debug("checked successfully")
	if *alreadyExists {
		n.Log.Error("Already exists")
		return errors.New("Already exists")
	}
	n.Log.Debug("doesnt already exist")

	n.Log.Debug("saving obj")
	err = n.db.SaveObject(sub)
	if err != nil {
		n.Log.Error(err.Error())
		return err
	}
	n.Log.Debug("saved")
	return nil
}
