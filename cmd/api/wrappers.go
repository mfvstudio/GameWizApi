package api

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/getkin/kin-openapi/openapi3filter"
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

func (app *Application) ValidateRequest(handler http.Handler) http.Handler {
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Validating request")
		route, pathParams, err := app.Config.RequestValidator.FindRoute(r)
		if err != nil {
			http.Error(w, "Route not found: "+err.Error(), http.StatusNotFound)
			return
		}
		input := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
			Options: &openapi3filter.Options{
				AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			},
		}
		if err = openapi3filter.ValidateRequest(r.Context(), input); err != nil {
			http.Error(w, "Request Validation Failed: "+err.Error(), http.StatusBadRequest)
			return
		}
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
