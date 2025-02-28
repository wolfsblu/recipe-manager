//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/handlers"
	"github.com/wolfsblu/go-chef/infra/job"
	"github.com/wolfsblu/go-chef/infra/routing"
	"github.com/wolfsblu/go-chef/infra/smtp"
	"github.com/wolfsblu/go-chef/infra/sqlite"
	"net/http"
)

func InitializeRecipeService() (*domain.RecipeService, error) {
	panic(wire.Build(
		domain.NewRecipeService,
		smtp.NewSMTPMailer,
		sqlite.NewSqliteStore,
		wire.Bind(new(domain.NotificationSender), new(*smtp.Mailer)),
		wire.Bind(new(domain.RecipeStore), new(*sqlite.Store)),
	))
}

func InitializeWebServer(service *domain.RecipeService) (*http.ServeMux, error) {
	panic(wire.Build(
		api.NewAPIServer,
		handlers.NewRecipeHandler,
		handlers.NewSecurityHandler,
		routing.NewServeMux,
		wire.Bind(new(api.Handler), new(*handlers.RecipeHandler)),
		wire.Bind(new(api.SecurityHandler), new(*handlers.SecurityHandler)),
	))
}

func InitializeScheduler(service *domain.RecipeService) *job.Scheduler {
	panic(wire.Build(job.NewScheduler))
}
