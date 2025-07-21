//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/handler"
	"github.com/wolfsblu/go-chef/infra/job"
	"github.com/wolfsblu/go-chef/infra/mapper"
	"github.com/wolfsblu/go-chef/infra/routing"
	"github.com/wolfsblu/go-chef/infra/smtp"
	"github.com/wolfsblu/go-chef/infra/sqlite"
	"github.com/wolfsblu/go-chef/infra/urlbuilder"
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
		mapper.NewRecipeMapper,
		urlbuilder.NewURLBuilder,
		handler.NewRecipeHandler,
		handler.NewSecurityHandler,
		handler.NewUploadHandler,
		routing.NewServeMux,
		wire.Bind(new(api.Handler), new(*handler.RecipeHandler)),
		wire.Bind(new(api.SecurityHandler), new(*handler.SecurityHandler)),
	))
}

func InitializeScheduler(service *domain.RecipeService) *job.Scheduler {
	panic(wire.Build(job.NewScheduler))
}
