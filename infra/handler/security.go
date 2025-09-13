package handler

import (
	"context"

	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/config"
)

type SecurityHandler struct {
	Users *domain.UserService
}

func NewSecurityHandler(service *domain.UserService) *SecurityHandler {
	return &SecurityHandler{
		Users: service,
	}
}

func (h *SecurityHandler) HandleCookieAuth(ctx context.Context, _ string, t api.CookieAuth) (context.Context, error) {
	userId, err := getUserFromSessionCookie(t.APIKey)
	if err != nil {
		return nil, domain.WrapError(domain.ErrAuthentication, err)
	}
	user, err := h.Users.GetUserById(ctx, userId)
	if err != nil {
		return nil, domain.WrapError(domain.ErrAuthentication, err)
	}
	return context.WithValue(ctx, config.CtxKeyUser, &user), nil
}
