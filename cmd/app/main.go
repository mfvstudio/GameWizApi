package main

import (
	"context"
	"fmt"
	"log"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/mfvstudio/gamewizapi/cmd/api"
	"github.com/mfvstudio/gamewizapi/internal/data"
)

func main() {
	fmt.Printf("Starting app...")
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("GameWizOpenApiSpec.yaml")
	if err != nil {
		panic(err)
	}
	if err = doc.Validate(context.Background()); err != nil {
		panic(err)
	}
	openApiRouter, err := gorillamux.NewRouter(doc)
	if err != nil {
		panic(err)
	}

	app := api.NewApp(&api.AppConfig{
		Port:             ":8080",
		Auth:             &data.FireAuthStore{},
		Data:             &data.FireDataStore{},
		RequestValidator: openApiRouter,
	})
	log.Fatal(app.Run(app.Mount()))
}
