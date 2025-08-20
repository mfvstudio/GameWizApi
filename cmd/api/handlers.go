package api

import (
	"encoding/json"
	"net/http"

	"github.com/mfvstudio/gamewizapi/cmd/gen"
)

func (Application) Health(w http.ResponseWriter, r *http.Request) {
	resp := gen.HealthCheck{
		Message: "Genki Desu!",
	}
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
