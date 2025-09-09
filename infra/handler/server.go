package handler

import (
	"github.com/wolfsblu/go-chef/domain"
)

type APIHandler struct {
	*RecipeHandler
	*UserHandler
}

func NewAPIHandler(recipes *domain.RecipeService, users *domain.UserService) *APIHandler {
	return &APIHandler{
		RecipeHandler: NewRecipeHandler(recipes),
		UserHandler:   NewUserHandler(users),
	}
}
