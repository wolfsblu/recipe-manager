package handler

import (
	"context"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
)

type contextKey string

const (
	ctxKeyUser = contextKey("User")
)

type SecurityHandler struct {
	Recipes *domain.RecipeService
}

func NewSecurityHandler(service *domain.RecipeService) *SecurityHandler {
	return &SecurityHandler{
		Recipes: service,
	}
}

func (h *SecurityHandler) HandleCookieAuth(ctx context.Context, _ string, t api.CookieAuth) (context.Context, error) {
	userId, err := getUserFromSessionCookie(t.APIKey)
	if err != nil {
		return nil, domain.WrapError(domain.ErrAuthentication, err)
	}
	user, err := h.Recipes.GetUserById(ctx, userId)
	if err != nil {
		return nil, domain.WrapError(domain.ErrAuthentication, err)
	}
	return context.WithValue(ctx, ctxKeyUser, &user), nil
}
