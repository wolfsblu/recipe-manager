package main

import (
	"github.com/wolfsblu/go-chef/infra/env"
	"github.com/wolfsblu/go-chef/infra/jobs"
	"github.com/wolfsblu/go-chef/infra/routing"
	"log"
	"net/http"
)

func main() {
	env.Load()

	recipeService, err := InitializeRecipeService()
	if err != nil {
		log.Fatalln("failed to initialize recipe service:", err)
	}

	jobs.StartScheduler(recipeService)
	defer jobs.QuitScheduler()

	apiServer, err := InitializeAPIServer()
	if err != nil {
		log.Fatalln("failed to initialize API server:", err)
	}
	mux := routing.NewServeMux(apiServer)

	host := env.MustGet("HOST")
	err = http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatalln("failed to start web server:", err)
	}
}
