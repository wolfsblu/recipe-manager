package handler

import "github.com/wolfsblu/go-chef/domain"

type Server struct {
	*RecipeHandler
	*UserHandler
}

func NewServer(recipes *domain.RecipeService, users *domain.UserService) *Server {
	return &Server{
		RecipeHandler: NewRecipeHandler(recipes),
		UserHandler:   NewUserHandler(users),
	}
}
