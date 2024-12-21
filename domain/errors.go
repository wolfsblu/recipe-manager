package domain

import (
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

var (
	ErrUnhandled                  = &Error{Code: 100, Message: "internal server error"}
	ErrStartingTransaction        = &Error{Code: 101, Message: "failed to establish transaction"}
	ErrCommittingTransaction      = &Error{Code: 102, Message: "failed to commit transaction"}
	ErrCreatingUser               = &Error{Code: 200, Message: "failed to create user"}
	ErrCreatingPasswordResetToken = &Error{Code: 201, Message: "failed to create password reset token"}
	ErrCreatingRegistrationToken  = &Error{Code: 202, Message: "failed to create user registration token"}
	ErrUserExists                 = &Error{Code: 203, Message: "user already exists"}
	ErrInvalidCredentials         = &Error{Code: 300, Message: "invalid credentials"}
	ErrUpdatingPassword           = &Error{Code: 301, Message: "failed to update password"}
	ErrAuthentication             = &Error{Code: 302, Message: "failed to authenticate user"}
	ErrDeletingPasswordResetToken = &Error{Code: 302, Message: "failed to remove password reset token"}
	ErrUserNotFound               = &Error{Code: 400, Message: "user was not found"}
	ErrPasswordResetTokenNotFound = &Error{Code: 401, Message: "password reset token was not found"}
	ErrRecipeNotFound             = &Error{Code: 402, Message: "recipe was not found"}
)

func (e *Error) Error() string {
	return e.Message
}

func WrapError(parent *Error, child error) error {
	return fmt.Errorf("%w: %w", parent, child)
}
