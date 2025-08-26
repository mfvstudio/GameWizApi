package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mfvstudio/gamewizapi/cmd/gen"
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
	res, err := app.Config.Data.GetGameSession(r, gameId)
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
