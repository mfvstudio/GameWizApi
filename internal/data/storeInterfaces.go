package data

import "net/http"

type AuthStore interface {
	CreateUserAccount(r *http.Request) error
}
