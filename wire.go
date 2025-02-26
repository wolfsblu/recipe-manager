//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/handlers"
	"github.com/wolfsblu/go-chef/infra/jobs"
	"github.com/wolfsblu/go-chef/infra/routing"
	"github.com/wolfsblu/go-chef/infra/smtp"
	"github.com/wolfsblu/go-chef/infra/sqlite"
	"net/http"
)

func InitializeRecipeService() (*domain.RecipeService, error) {
	panic(wire.Build(
		smtp.Set,
		sqlite.Set,
		domain.NewRecipeService,
	))
}

func InitializeWebServer(service *domain.RecipeService) (*http.ServeMux, error) {
	panic(wire.Build(
		handlers.Set,
		api.NewAPIServer,
		routing.NewServeMux,
	))
}

func InitializeScheduler(service *domain.RecipeService) *jobs.Scheduler {
	panic(wire.Build(jobs.NewScheduler))
}
