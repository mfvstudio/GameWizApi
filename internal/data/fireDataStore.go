package data

import (
	"context"
	"errors"
	"log"
	"net/http"
	"slices"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/mfvstudio/gamewizapi/cmd/gen"
	"github.com/mfvstudio/gamewizapi/internal/env"
)

var dataOnce sync.Once

type FireDataStore struct {
	dataClient *firestore.Client
}

func (fs *FireDataStore) client() *firestore.Client {
	dataOnce.Do(func() {
		str, err := env.GetString("GOOGLE_CLOUD_PROJECT")
		var config *firebase.Config
		if err == nil {
			config = &firebase.Config{ProjectID: str}
		}
		app, err := firebase.NewApp(context.Background(), config)
		if err != nil {
			log.Fatalf("Error initializing data app: %v\n", err)
		}
		store, err := app.Firestore(context.Background())
		if err != nil {
			log.Fatalf("Unable to init firestore client: %v\n", err)
		}
		fs.dataClient = store
	})
	return fs.dataClient
}

func (fs *FireDataStore) GetGameSession(r *http.Request, gameId string) (*gen.Session, error) {
	doc, err := fs.client().Collection("sessions").Doc(gameId).Get(r.Context())
	if err != nil {
		log.Printf("No results returned: %v", err)
		return nil, err
	}
	var res gen.Session
	doc.DataTo(&res)
	return &res, nil
}

func (fs *FireDataStore) PutGameSession(s *gen.Session) error {
	result, err := fs.client().Collection("sessions").Doc(s.SessionId).Set(context.Background(), s)
	if err != nil {
		log.Printf("Error while inserting new session to data store: %v", err)
		return err
	}
	log.Printf("New session creation success: %v", result)
	return nil
}

func (fs *FireDataStore) JoinSession(s *gen.JoinSession) (*gen.Session, error) {
	//try to find session
	// if it does not exist, return error
	result, err := fs.client().Collection("inviteCodes").Doc(s.InviteCode).Get(context.Background())
	if err != nil {
		return nil, err
	}
	sessionStruct := struct {
		SessionId string `json:"sessionId"`
		PlayerUID string `json:"playerUID"`
	}{}
	result.DataTo(&sessionStruct)
	session, err := fs.client().Collection("sessions").Doc(sessionStruct.SessionId).Get(context.Background())
	if err != nil {
		return nil, err
	}
	var final gen.Session
	session.DataTo(&final)
	if len(final.Players) == final.MaxPlayerCount {
		return nil, errors.New("MaxCapacityReached")
	}
	if slices.Contains(final.Players, sessionStruct.PlayerUID) {
		return nil, errors.New("InvalidRequest")
	}
	final.Players = append(final.Players, sessionStruct.PlayerUID)
	update := []firestore.Update{
		{Path: "players", Value: final.Players},
	}
	if len(final.Players) == final.MaxPlayerCount {
		newSession := gen.INPROGRESS
		update = append(update, firestore.Update{Path: "status", Value: newSession})
	}
	_, err = fs.client().Collection("sessions").Doc(final.SessionId).Update(context.Background(), update)
	if err != nil {
		return nil, err
	}
	return &final, nil
}
