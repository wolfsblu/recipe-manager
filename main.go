package main

import (
	"log"
	"net/http"

	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/env"
	"github.com/wolfsblu/recipe-manager/infra/handler"
	"github.com/wolfsblu/recipe-manager/infra/job"
	"github.com/wolfsblu/recipe-manager/infra/routing"
	"github.com/wolfsblu/recipe-manager/infra/smtp"
	"github.com/wolfsblu/recipe-manager/infra/sqlite"
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
	shoppingService := domain.NewShoppingService(sqliteStore)

	securityHandler := handler.NewSecurityHandler(userService)
	apiHandler := handler.NewAPIHandler(recipeService, userService, shoppingService)
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
