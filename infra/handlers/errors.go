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

var errorStatusCodeMap = map[*domain.Error]int{
	domain.ErrAuthentication:             http.StatusForbidden,
	domain.ErrCommittingTransaction:      http.StatusInternalServerError,
	domain.ErrCreatingPasswordResetToken: http.StatusInternalServerError,
	domain.ErrCreatingRegistrationToken:  http.StatusInternalServerError,
	domain.ErrCreatingUser:               http.StatusInternalServerError,
	domain.ErrDeletingPasswordResetToken: http.StatusInternalServerError,
	domain.ErrDeletingRegistration:       http.StatusInternalServerError,
	domain.ErrInvalidCredentials:         http.StatusForbidden,
	domain.ErrPasswordResetTokenNotFound: http.StatusForbidden,
	domain.ErrRecipeNotFound:             http.StatusNotFound,
	domain.ErrRegistrationNotFound:       http.StatusNotFound,
	domain.ErrStartingTransaction:        http.StatusInternalServerError,
	domain.ErrUnhandled:                  http.StatusInternalServerError,
	domain.ErrUpdatingPassword:           http.StatusInternalServerError,
	domain.ErrUpdatingUser:               http.StatusInternalServerError,
	domain.ErrUserExists:                 http.StatusBadRequest,
	domain.ErrUserNotFound:               http.StatusNotFound,
}

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
	if val, ok := errorStatusCodeMap[err]; ok {
		return val
	}
	return http.StatusInternalServerError
}
