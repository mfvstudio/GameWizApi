package api

import (
	"net/http"

	"slices"

	"github.com/mfvstudio/gamewizapi/cmd/gen"
)

func IsAdmin(r *http.Request) bool {
	claims := r.Context().Value("userToken").(*AuthToken)
	return claims.Context.Roles == "admin"
}

func IsUser(r *http.Request) bool {
	claims := r.Context().Value("userToken").(*AuthToken)
	return claims.Context.Roles == "user"
}

func IsPlayerInSession(s *gen.Session, r *http.Request) bool {
	claims := r.Context().Value("userToken").(*AuthToken)
	return slices.Contains(*s.Players, claims.UID)
}
