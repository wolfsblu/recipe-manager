package main

import (
	"log"
	"net/http"

	"github.com/wolfsblu/go-chef/infra/env"
)

func main() {
	env.Load()
	recipeService, err := InitializeRecipeService()
	if err != nil {
		log.Fatal("failed to initialize recipe service: ", err)
	}

	scheduler := InitializeScheduler(recipeService)
	defer scheduler.Quit()

	mux, err := InitializeWebServer(recipeService)
	if err != nil {
		log.Fatal("failed to initialize web server: ", err)
	}

	host := env.MustGet("HOST")
	err = http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatalln("failed to start web server: ", err)
	}
}
