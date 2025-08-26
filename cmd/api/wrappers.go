package api

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
)

type AuthToken struct {
	UID     string  `json:"user_id"`
	Context Context `json:"context"`
}

type Context struct {
	Roles string `json:"roles"`
}

func SetBasicHeaders(handler http.Handler) http.Handler {
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		handler.ServeHTTP(w, r)
	}))
}

func (app *Application) AuthenticateToken(handler http.Handler) http.Handler {
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Authenticating request")
		token, err := app.Config.Auth.Authenticate(r)
		if err != nil {
			if err.Error() == "MissingAuthHeader" {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		log.Printf("Returned token: %v", token)
		tokenStruct := tokenClaimsToAuthStruct(token)
		ctx := context.WithValue(r.Context(), "userToken", tokenStruct)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}))
}

func tokenClaimsToAuthStruct(u *auth.Token) *AuthToken {
	t := AuthToken{
		UID: u.UID,
		Context: Context{
			Roles: "",
		},
	}
	context := u.Claims["context"].(map[string]interface{})
	if roles, ok := context["roles"].(string); ok {
		t.Context.Roles = roles
	}
	return &t
}
