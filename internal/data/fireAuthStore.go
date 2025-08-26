package data

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/mfvstudio/gamewizapi/cmd/gen"
	"github.com/mfvstudio/gamewizapi/internal/env"
	"github.com/mfvstudio/gamewizapi/internal/helpers"
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

func (fas *FireAuthStore) CreateUserAccount(r *http.Request) (*string, error) {
	var form *gen.CreateAccount
	json.NewDecoder(r.Body).Decode(&form)
	accepted, err := helpers.IsSignUpFormAccepted(form)
	if err != nil {
		return nil, err
	}
	if !accepted {
		log.Printf("Form failed auth checks")
		return nil, errors.New("FormNotAccepted")
	}
	usrReq := (&auth.UserToCreate{}).
		Email(string(form.Email)).
		EmailVerified(false).
		Password(form.Password).
		DisplayName(form.Username).
		Disabled(false)
	log.Printf("User request is: %v", usrReq)
	record, err := fas.client().CreateUser(r.Context(), usrReq)
	if err != nil {
		return nil, err
	}
	log.Printf("User signup Result: %v", record)
	verify, err := fas.client().EmailVerificationLink(r.Context(), string(form.Email))
	if err != nil {
		log.Printf("Error in creating email verification link")
		return nil, err
	}
	return &verify, nil
}

func (fas *FireAuthStore) Authenticate(r *http.Request) (*auth.Token, error) {
	jwt := r.Header.Get("Authorization")
	if jwt == "" {
		return nil, errors.New("MissingAuthHeader")
	}
	token, err := fas.client().VerifyIDToken(r.Context(), jwt)
	if err != nil {
		log.Printf("Error while authenticating jwt through authStore: %v", err)
		return nil, errors.New("InternalServerError")
	}
	return token, nil
}
