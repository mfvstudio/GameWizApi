package main

import (
	"fmt"
	"log"

	"github.com/mfvstudio/gamewizapi/cmd/api"
	"github.com/mfvstudio/gamewizapi/internal/data"
)

func main() {
	fmt.Printf("Starting app...")
	app := api.NewApp(&api.AppConfig{
		Port: ":8080",
		Auth: &data.FireAuthStore{},
	})
	log.Fatal(app.Run(app.Mount()))
}
