package api

import (
	"log"
	"net/http"
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

func (app *Application) Run() error {
	port := app.Config.Port
	server := &http.Server{
		Addr:    port,
		Handler: nil, //TODO: add generated Mux handler
	}
	log.Printf("\nServer Listening on port %s\n", port)
	log.Printf(`GameWizApi`)
	return server.ListenAndServe()
}
