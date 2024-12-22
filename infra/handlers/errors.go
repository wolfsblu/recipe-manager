package handlers

import (
	"context"
	"encoding/json"
	"errors"
	ogenhttp "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"net/http"
	"strings"
)

const (
	JsonIndentationChar   = " "
	JsonIndentationPrefix = ""
	JsonIndentationWidth  = 2
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
		StatusCode: mapDomainErrorToStatusCode(domainErr),
		Response: api.Error{
			Message: domainErr.Message,
		},
	}
}

func (h *RecipeHandler) CustomErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	var (
		code    int
		message string
		ogenErr ogenerrors.Error
	)

	switch {
	case errors.Is(err, ogenhttp.ErrNotImplemented):
		code = http.StatusNotImplemented
		message = http.StatusText(http.StatusNotImplemented)
	case errors.As(err, &ogenErr):
		code = ogenErr.Code()
		message = ogenErr.Error()
	default:
		code = http.StatusInternalServerError
		message = err.Error()
	}

	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.SetIndent(JsonIndentationPrefix, strings.Repeat(JsonIndentationChar, JsonIndentationWidth))
	_ = encoder.Encode(api.Error{
		Message: message,
	})
}

func mapDomainErrorToStatusCode(err *domain.Error) int {
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
