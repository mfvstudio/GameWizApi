package main

import (
	"fmt"
	"log"

	"github.com/mfvstudio/gamewizapi/cmd/api"
)

func main() {
	fmt.Printf("Starting app...")
	app := api.NewApp(&api.AppConfig{
		Port: ":8080",
	})
	log.Fatal(app.Run())
}
