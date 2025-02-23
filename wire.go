//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/handlers"
	"github.com/wolfsblu/go-chef/infra/smtp"
	"github.com/wolfsblu/go-chef/infra/sqlite"
)

var recipeServiceSet = wire.NewSet(
	smtp.Set,
	sqlite.Set,
	domain.NewRecipeService,
)

func InitializeRecipeService() (*domain.RecipeService, error) {
	wire.Build(recipeServiceSet)
	return &domain.RecipeService{}, nil
}

func InitializeAPIServer() (*api.Server, error) {
	wire.Build(
		handlers.Set,
		recipeServiceSet,
		api.NewAPIServer,
	)
	return &api.Server{}, nil
}
