package data

import (
	"context"
	"log"
	"net/http"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/mfvstudio/gamewizapi/internal/env"
)

var authOnce sync.Once

type FireAuthStore struct {
	authClient *auth.Client
}

func (fas *FireAuthStore) client() *auth.Client {
	authOnce.Do(func() {
		str, err := env.GetString("GOOGLE_CLOUD_PROJECT")
		var config *firebase.Config
		if err == nil {
			config = &firebase.Config{ProjectID: str}
		}
		app, err := firebase.NewApp(context.Background(), config)
		if err != nil {
			log.Fatalf("Error initializing data app: %v\n", err)
		}
		store, err := app.Auth(context.Background())
		if err != nil {
			log.Fatalf("Unable to init firestore client: %v\n", err)
		}
		fas.authClient = store
	})
	return fas.authClient
}

func (fas *FireAuthStore) CreateUserAccount(r *http.Request) error {
	//TODO: ;)
	return nil
}
