package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
)

var errorStatusCodeMap = map[*domain.Error]int{
	domain.ErrAuthentication:             http.StatusUnauthorized,
	domain.ErrAuthorization:              http.StatusForbidden,
	domain.ErrCommittingTransaction:      http.StatusInternalServerError,
	domain.ErrCreatingPasswordResetToken: http.StatusInternalServerError,
	domain.ErrCreatingRegistrationToken:  http.StatusInternalServerError,
	domain.ErrCreatingUser:               http.StatusInternalServerError,
	domain.ErrDeletingPasswordResetToken: http.StatusInternalServerError,
	domain.ErrDeletingRegistration:       http.StatusInternalServerError,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrPagination:                 http.StatusBadRequest,
	domain.ErrPasswordResetTokenNotFound: http.StatusUnauthorized,
	domain.ErrRecipeNotFound:             http.StatusNotFound,
	domain.ErrRegistrationNotFound:       http.StatusNotFound,
	domain.ErrStartingTransaction:        http.StatusInternalServerError,
	domain.ErrUnconfirmedUser:            http.StatusForbidden,
	domain.ErrUnhandled:                  http.StatusInternalServerError,
	domain.ErrUpdatingPassword:           http.StatusInternalServerError,
	domain.ErrUpdatingUser:               http.StatusInternalServerError,
	domain.ErrUserExists:                 http.StatusBadRequest,
	domain.ErrUserNotFound:               http.StatusNotFound,
}

func (h *RecipeHandler) NewError(_ context.Context, err error) (r *api.ErrorStatusCode) {
	var domainErr = domain.ErrUnhandled
	var securityError *ogenerrors.SecurityError

	if errors.As(err, &securityError) {
		domainErr = domain.ErrAuthentication
	} else if !errors.As(err, &domainErr) {
		log.Println(fmt.Errorf("unhandled error: %w", err))
		domainErr = domain.ErrUnhandled
	}

	return &api.ErrorStatusCode{
		StatusCode: mapDomainErrorToStatusCode(domainErr),
		Response: api.Error{
			Message: domainErr.Message,
		},
	}
}

func mapDomainErrorToStatusCode(err *domain.Error) int {
	if val, ok := errorStatusCodeMap[err]; ok {
		return val
	}
	return http.StatusInternalServerError
}
