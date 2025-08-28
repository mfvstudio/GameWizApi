package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mfvstudio/gamewizapi/cmd/gen"
	"github.com/mfvstudio/gamewizapi/internal/helpers"
	"github.com/mfvstudio/gamewizapi/internal/wizErrors"
)

func (app *Application) Health(w http.ResponseWriter, r *http.Request) {
	if !IsAdmin(r) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	resp := gen.HealthCheck{
		Message: "Genki Desu!",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (app *Application) CreateAccount(w http.ResponseWriter, r *http.Request) {
	verifyEmail, err := app.Config.Auth.CreateUserAccount(r)
	if err != nil {
		log.Printf("Error creating new user account: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := gen.ApiResponse{
		Message: verifyEmail,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (app *Application) GetSession(w http.ResponseWriter, r *http.Request, gameId string) {
	//Check if caller has correct permissions
	res, err := app.Config.Data.GetGameSession(gameId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !IsPlayerInSession(res, r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (app *Application) CreateSession(w http.ResponseWriter, r *http.Request) {
	var request gen.CreateSession
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session, err := helpers.CreateSessionFromSessionRequest(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	inviteResp, err := app.Config.Data.PutGameSession(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&inviteResp)
}

func (app *Application) JoinSession(w http.ResponseWriter, r *http.Request) {
	var s gen.JoinSession
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session, err := app.Config.Data.JoinSession(&s)
	if err != nil {
		switch err {
		case wizErrors.ResourceNotFound:
			log.Printf("Invite or Session id not found: %v", err)
			w.WriteHeader(http.StatusNotFound)
		case wizErrors.MaxCapacityReached:
			w.WriteHeader(http.StatusNotAcceptable)
		case wizErrors.UserAlreadyInSession:
			w.WriteHeader(http.StatusNotAcceptable)
		default:
			log.Printf("Server error while joining session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func (app *Application) UpdateSession(w http.ResponseWriter, r *http.Request, gameId string) {
	var s gen.UpdateSession
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := app.Config.Data.UpdateSession(&s, gameId); err != nil {
		log.Printf("Error while updating a session: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
