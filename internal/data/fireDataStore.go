package data

import (
	"context"
	"log"
	"slices"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/mfvstudio/gamewizapi/cmd/gen"
	"github.com/mfvstudio/gamewizapi/internal/env"
	"github.com/mfvstudio/gamewizapi/internal/helpers"
	"github.com/mfvstudio/gamewizapi/internal/wizErrors"
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

func (fs *FireDataStore) GetGameSession(gameId string) (*gen.Session, error) {
	doc, err := fs.client().Collection("sessions").Doc(gameId).Get(context.Background())
	if err != nil {
		log.Printf("No results returned: %v", err)
		return nil, err
	}
	var res gen.Session
	doc.DataTo(&res)
	return &res, nil
}

func (fs *FireDataStore) PutGameSession(s *gen.Session) (*gen.InviteCodeResponse, error) {
	result, err := fs.client().Collection("sessions").Doc(s.SessionId).Set(context.Background(), s)
	if err != nil {
		log.Printf("Error while inserting new session to data store: %v", err)
		return nil, err
	}
	log.Printf("New session creation success: %v", result)
	//TODO: Should we retry to recreate invite codes? the chances of a collision are 26^6. Which is highly unlikely.
	//we can add a TTL to the invite codes so we can delete them and minimize this.
	inviteCode := helpers.GenerateSessionInviteCode()

	invite := gen.InviteCodeResponse{
		InviteCode: inviteCode,
	}
	_, err = fs.client().Collection("inviteCodes").Doc(inviteCode).Set(context.Background(), struct {
		SessionId string `json:"sessionId"`
	}{SessionId: s.SessionId})
	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (fs *FireDataStore) JoinSession(s *gen.JoinSession) (*gen.Session, error) {
	result, err := fs.client().Collection("inviteCodes").Doc(s.InviteCode).Get(context.Background())
	if err != nil {
		return nil, wizErrors.ResourceNotFound
	}
	sessionStruct := struct {
		SessionId string `json:"sessionId"`
	}{}
	result.DataTo(&sessionStruct)
	session, err := fs.client().Collection("sessions").Doc(sessionStruct.SessionId).Get(context.Background())
	if err != nil {
		return nil, wizErrors.ResourceNotFound
	}
	var final gen.Session
	session.DataTo(&final)
	if len(final.Players) == final.MaxPlayerCount {
		return nil, wizErrors.MaxCapacityReached
	}
	if slices.Contains(final.Players, s.PlayerUID) {
		return nil, wizErrors.UserAlreadyInSession
	}
	final.Players = append(final.Players, s.PlayerUID)
	update := []firestore.Update{
		{Path: "Players", Value: final.Players},
	}
	if len(final.Players) == final.MaxPlayerCount {
		newSession := gen.INPROGRESS
		update = append(update, firestore.Update{Path: "Status", Value: newSession})
		final.Status = newSession
	}
	_, err = fs.client().Collection("sessions").Doc(final.SessionId).Update(context.Background(), update)
	if err != nil {
		return nil, err
	}
	return &final, nil
}

func (fs *FireDataStore) UpdateSession(s *gen.UpdateSession, gameId string) error {
	update := []firestore.Update{
		{Path: "GameState", Value: s.GameState},
	}
	if s.MetaData != nil {
		update = append(update, firestore.Update{Path: "MetaData", Value: s.MetaData})
	}
	_, err := fs.client().Collection("sessions").Doc(gameId).Update(context.Background(), update)
	if err != nil {
		return err
	}
	return nil
}
