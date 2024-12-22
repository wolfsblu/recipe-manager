package domain

type Error struct {
	Inner   error
	Message string
}

var (
	ErrAuthentication             = &Error{Message: "failed to authenticate user"}
	ErrCommittingTransaction      = &Error{Message: "failed to commit transaction"}
	ErrCreatingPasswordResetToken = &Error{Message: "failed to create password reset token"}
	ErrCreatingRegistrationToken  = &Error{Message: "failed to create user registration token"}
	ErrCreatingUser               = &Error{Message: "failed to create user"}
	ErrDeletingPasswordResetToken = &Error{Message: "failed to remove password reset token"}
	ErrDeletingRegistration       = &Error{Message: "failed to complete user registration"}
	ErrInvalidCredentials         = &Error{Message: "invalid credentials"}
	ErrPasswordResetTokenNotFound = &Error{Message: "password reset token was not found"}
	ErrRecipeNotFound             = &Error{Message: "recipe was not found"}
	ErrRegistrationNotFound       = &Error{Message: "user registration was not found"}
	ErrStartingTransaction        = &Error{Message: "failed to establish transaction"}
	ErrUnconfirmedUser            = &Error{Message: "the requested user is not confirmed"}
	ErrUnhandled                  = &Error{Message: "internal server error"}
	ErrUpdatingPassword           = &Error{Message: "failed to update password"}
	ErrUpdatingUser               = &Error{Message: "failed to update user"}
	ErrUserExists                 = &Error{Message: "user already exists"}
	ErrUserNotFound               = &Error{Message: "user was not found"}
)

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Inner
}

func WrapError(err *Error, inner error) error {
	return &Error{
		Inner:   inner,
		Message: err.Error(),
	}
}
