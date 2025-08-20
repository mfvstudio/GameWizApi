package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mfvstudio/gamewizapi/cmd/gen"
)

type Application struct {
	Config AppConfig
}

type AppConfig struct {
	Port string
}

func NewApp(config *AppConfig) Application {
	return Application{
		Config: *config,
	}
}

func (Application) Mount() *chi.Mux {
	r := chi.NewMux()
	return r
}

func (app *Application) Run(mux *chi.Mux) error {
	h := gen.HandlerFromMuxWithBaseURL(app, mux, "/v0") //TODO: read api version from somewhere
	port := app.Config.Port
	server := &http.Server{
		Addr:    port,
		Handler: h,
	}
	log.Printf("\nServer Listening on port %s\n", port)
	log.Printf(`
 ________  ________  _____ ______   _______   ___       __   ___  ________     
|\   ____\|\   __  \|\   _ \  _   \|\  ___ \ |\  \     |\  \|\  \|\_____  \    
\ \  \___|\ \  \|\  \ \  \\\__\ \  \ \   __/|\ \  \    \ \  \ \  \\|___/  /|   
 \ \  \  __\ \   __  \ \  \\|__| \  \ \  \_|/_\ \  \  __\ \  \ \  \   /  / /   
  \ \  \|\  \ \  \ \  \ \  \    \ \  \ \  \_|\ \ \  \|\__\_\  \ \  \ /  /_/__  
   \ \_______\ \__\ \__\ \__\    \ \__\ \_______\ \____________\ \__\\________\
    \|_______|\|__|\|__|\|__|     \|__|\|_______|\|____________|\|__|\|_______|
	
`)
	return server.ListenAndServe()
}
