//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/handlers"
	"github.com/wolfsblu/go-chef/infra/jobs"
	"github.com/wolfsblu/go-chef/infra/smtp"
	"github.com/wolfsblu/go-chef/infra/sqlite"
)

var recipeServiceSet = wire.NewSet(
	smtp.Set,
	sqlite.Set,
	domain.NewRecipeService,
)

func InitializeAPIServer() (*api.Server, error) {
	wire.Build(
		handlers.Set,
		recipeServiceSet,
		api.NewAPIServer,
	)
	return &api.Server{}, nil
}

func InitializeScheduler() (*jobs.Scheduler, error) {
	wire.Build(recipeServiceSet, jobs.NewScheduler)
	return &jobs.Scheduler{}, nil
}
