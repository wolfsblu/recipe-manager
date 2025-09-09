package handler

import (
	"context"

	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/domain/security"
	"github.com/wolfsblu/go-chef/infra/config"
)

type UserHandler struct {
	Users *domain.UserService
}

func NewUserHandler(service *domain.UserService) *UserHandler {
	return &UserHandler{
		Users: service,
	}
}

func (h *UserHandler) ConfirmUser(ctx context.Context, req *api.Token) error {
	return h.Users.ConfirmUserByToken(ctx, req.Token)
}

func (h *UserHandler) GetUserProfile(ctx context.Context) (*api.ReadUser, error) {
	user := ctx.Value(config.CtxKeyUser).(*domain.User)
	return &api.ReadUser{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *api.Credentials) (r *api.AuthenticatedUserHeaders, _ error) {
	user, err := h.Users.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = h.Users.VerifyPassword(user, req.Password)
	if err != nil {
		return nil, err
	} else if !user.Confirmed {
		return nil, domain.ErrUnconfirmedUser
	}

	cookie, err := createSessionCookie(user.ID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrAuthentication, err)
	}
	return &api.AuthenticatedUserHeaders{
		SetCookie: api.OptString{
			Set:   true,
			Value: cookie,
		},
		Response: api.ReadUser{
			ID:    user.ID,
			Email: user.Email,
		},
	}, nil
}

func (h *UserHandler) Logout(_ context.Context) (*api.LogoutOK, error) {
	cookie := expireSessionCookie()
	return &api.LogoutOK{
		SetCookie: api.OptString{
			Set:   true,
			Value: cookie,
		},
	}, nil
}

func (h *UserHandler) Register(ctx context.Context, c *api.Credentials) error {
	hash, err := security.CreateHash(c.Password, security.DefaultHashParams)
	if err != nil {
		return domain.WrapError(domain.ErrCreatingUser, err)
	}
	return h.Users.RegisterUser(ctx, domain.Credentials{
		Email:        c.Email,
		PasswordHash: hash,
	})
}

func (h *UserHandler) ResetPassword(ctx context.Context, req *api.ResetPasswordReq) error {
	return h.Users.ResetPasswordByEmail(ctx, req.Email)
}

func (h *UserHandler) UpdatePassword(ctx context.Context, req *api.PasswordReset) error {
	hashedPassword, err := security.CreateHash(req.Password, security.DefaultHashParams)
	if err != nil {
		return err
	}
	return h.Users.UpdatePasswordByToken(ctx, req.Token, hashedPassword)
}
