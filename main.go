package main

import (
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/env"
	"github.com/wolfsblu/go-chef/infra/handlers"
	"github.com/wolfsblu/go-chef/infra/jobs"
	"github.com/wolfsblu/go-chef/infra/routing"
	"github.com/wolfsblu/go-chef/infra/smtp"
	"github.com/wolfsblu/go-chef/infra/sqlite"
	"log"
	"net/http"
)

func main() {
	env.Load()

	notifier := smtp.NewSMTPMailer()
	store, err := sqlite.NewSqliteStore()
	if err != nil {
		log.Fatalln("failed to initialize sqlite connection:", err)
	}
	err = store.Migrate()
	if err != nil {
		log.Fatalln("failed to apply database migrations:", err)
	}

	recipeService := domain.NewRecipeService(notifier, store)
	recipeHandler := handlers.NewRecipeHandler(recipeService)
	securityHandler := handlers.NewSecurityHandler(recipeService)
	errorHandler := recipeHandler.CustomErrorHandler
	apiServer, err := api.NewServer(recipeHandler, securityHandler, api.WithErrorHandler(errorHandler))
	if err != nil {
		log.Fatalln("failed to start api server:", err)
	}

	jobs.StartScheduler(recipeService)
	defer jobs.QuitScheduler()

	host := env.MustGet("HOST")
	mux := routing.NewServeMux(apiServer)
	err = http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatalln("failed to start web server:", err)
	}
}
