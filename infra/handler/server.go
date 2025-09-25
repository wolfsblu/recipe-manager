package handler

import (
	"github.com/wolfsblu/recipe-manager/domain"
)

type APIHandler struct {
	*RecipeHandler
	*UserHandler
	*ShoppingHandler
}

func NewAPIHandler(recipes *domain.RecipeService, users *domain.UserService, shopping *domain.ShoppingService) *APIHandler {
	return &APIHandler{
		RecipeHandler:   NewRecipeHandler(recipes),
		UserHandler:     NewUserHandler(users),
		ShoppingHandler: NewShoppingHandler(shopping),
	}
}
