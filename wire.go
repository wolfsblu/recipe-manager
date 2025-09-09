//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/handler"
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

func InitializeUserService() (*domain.UserService, error) {
	panic(wire.Build(
		domain.NewUserService,
		smtp.NewSMTPMailer,
		sqlite.NewSqliteStore,
		wire.Bind(new(domain.NotificationSender), new(*smtp.Mailer)),
		wire.Bind(new(domain.UserStore), new(*sqlite.Store)),
	))
}

func InitializeWebServer(recipeService *domain.RecipeService, userService *domain.UserService) (*http.ServeMux, error) {
	panic(wire.Build(
		api.NewAPIServer,
		handler.NewServer,
		handler.NewSecurityHandler,
		handler.NewUploadHandler,
		routing.NewServeMux,
		wire.Bind(new(api.Handler), new(*handler.Server)),
		wire.Bind(new(api.SecurityHandler), new(*handler.SecurityHandler)),
	))
}

func InitializeScheduler(userService *domain.UserService) *job.Scheduler {
	panic(wire.Build(job.NewScheduler))
}
