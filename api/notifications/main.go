package notifications

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/will-lol/personal_website/subscriptions"

	"encoding/json"
	"errors"
	"fmt"
	webpush "github.com/SherClockHolmes/webpush-go"
	"io"
	"log"
	"os"
)

var s *subscriptions.SubscriptionClient

func Router(r chi.Router) {
	r.Post("/subscribe", subscribeHandler)
	r.Get("/notify", notifyHandler)
	r.Get("/publickey", publicKeyHandler)
}

func getSecrets() (publicKey *string, privateKey *string, err error) {
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

	//The response is a JSON object that needs unmarshalling
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

func getClient() (*subscriptions.SubscriptionClient, error) {
	if s != nil {
		return s, nil
	}
	var err error
	s, err = subscriptions.NewSubscriptionClient()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	s, err := getClient()
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	subscribers, err := s.GetSubscribers()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	vapidPub, vapidPriv, err := getSecrets()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	for _, element := range *subscribers {
		go notifySubscriber(element, *vapidPub, *vapidPriv)
	}
	w.WriteHeader(200)
}

func notifySubscriber(sub webpush.Subscription, vapidPub string, vapidPriv string) error {
	resp, err := webpush.SendNotification([]byte("Notification received"), &sub, &webpush.Options{
		VAPIDPublicKey:  vapidPub,
		VAPIDPrivateKey: vapidPriv,
		TTL:             64,
		Subscriber:      "will.bradshaw50@gmail.com",
	})
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.StatusCode == 410 {
		s.Unsubscribe(sub.Endpoint)
	}
	return nil
}

func publicKeyHandler(w http.ResponseWriter, r *http.Request) {
	pub, _, err := getSecrets()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	w.WriteHeader(200)
	w.Write([]byte(*pub))
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	s, err := getClient()
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	subscriptionBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	var subscription webpush.Subscription
	json.Unmarshal(subscriptionBytes, &subscription)

	_, err = url.ParseRequestURI(subscription.Endpoint)
	if len(subscription.Endpoint) < 1 || err != nil {
		log.Fatalln(err)
		w.WriteHeader(400)
		w.Write([]byte("Improper URL"))
		return
	}

	alreadyExists, err := s.DoesSubscriptionExist(subscription.Endpoint)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		w.Write([]byte("Couldn't check the subscription"))
		return
	}
	if *alreadyExists {
		w.WriteHeader(409)
		w.Write([]byte("Subscription already exists"))
		return
	}

	err = s.SaveSubscription(subscription)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Success"))
}
