package handlers

import (
	"context"
	"errors"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
)

func (h *RecipeHandler) NewError(_ context.Context, err error) (r *api.ErrorStatusCode) {
	var domainErr *domain.Error
	var securityError *ogenerrors.SecurityError

	if errors.As(err, &securityError) {
		domainErr = domain.ErrAuthentication
	} else if !errors.As(err, &domainErr) {
		domainErr = domain.ErrUnhandled
	}

	return &api.ErrorStatusCode{
		StatusCode: mapToStatusCode(domainErr),
		Response: api.Error{
			Message: domainErr.Message,
		},
	}
}

func mapToStatusCode(err *domain.Error) int {
	switch {
	case errors.Is(err, domain.ErrUserExists):
		return 400
	case errors.Is(err, domain.ErrAuthentication):
		fallthrough
	case errors.Is(err, domain.ErrPasswordResetTokenNotFound):
		fallthrough
	case errors.Is(err, domain.ErrInvalidCredentials):
		return 403
	case errors.Is(err, domain.ErrRecipeNotFound):
		fallthrough
	case errors.Is(err, domain.ErrRegistrationNotFound):
		fallthrough
	case errors.Is(err, domain.ErrUserNotFound):
		return 404
	case errors.Is(err, domain.ErrCommittingTransaction):
		fallthrough
	case errors.Is(err, domain.ErrCreatingUser):
		fallthrough
	case errors.Is(err, domain.ErrCreatingPasswordResetToken):
		fallthrough
	case errors.Is(err, domain.ErrCreatingRegistrationToken):
		fallthrough
	case errors.Is(err, domain.ErrDeletingPasswordResetToken):
		fallthrough
	case errors.Is(err, domain.ErrDeletingRegistration):
		fallthrough
	case errors.Is(err, domain.ErrStartingTransaction):
		fallthrough
	case errors.Is(err, domain.ErrUpdatingUser):
		fallthrough
	case errors.Is(err, domain.ErrUnhandled):
		fallthrough
	default:
		return 500
	}
}
