package data

import (
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/mfvstudio/gamewizapi/cmd/gen"
)

type AuthStore interface {
	CreateUserAccount(r *http.Request) (*string, error)
	Authenticate(r *http.Request) (*auth.Token, error)
}

type DataStore interface {
	GetGameSession(r *http.Request, gameId string) (*gen.Session, error)
	PutGameSession(s *gen.Session) error
	JoinSession(s *gen.JoinSession) (*gen.Session, error)
}
