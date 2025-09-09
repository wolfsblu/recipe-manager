package main

import (
	"log"
	"net/http"

	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/env"
	"github.com/wolfsblu/go-chef/infra/handler"
	"github.com/wolfsblu/go-chef/infra/job"
	"github.com/wolfsblu/go-chef/infra/routing"
	"github.com/wolfsblu/go-chef/infra/smtp"
	"github.com/wolfsblu/go-chef/infra/sqlite"
)

func main() {
	env.Load()

	mailer := smtp.NewSMTPMailer()
	sqliteStore, err := sqlite.NewSqliteStore()
	if err != nil {
		log.Fatal("failed to initialize sqlite store: ", err)
	}

	recipeService := domain.NewRecipeService(mailer, sqliteStore)
	userService := domain.NewUserService(mailer, sqliteStore)

	securityHandler := handler.NewSecurityHandler(userService)
	apiHandler := handler.NewAPIHandler(recipeService, userService)
	uploadHandler, err := handler.NewUploadHandler(userService)
	if err != nil {
		log.Fatal("failed to initialize upload handler: ", err)
	}
	apiServer, err := api.NewAPIServer(apiHandler, securityHandler)
	if err != nil {
		log.Fatal("failed to initialize API server: ", err)
	}

	mux := routing.NewServeMux(apiServer, uploadHandler)
	scheduler := job.NewScheduler(userService)
	defer scheduler.Quit()

	host := env.MustGet("HOST")
	err = http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatalln("failed to start web server: ", err)
	}
}
