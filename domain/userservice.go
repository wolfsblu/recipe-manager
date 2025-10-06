package domain

import (
	"context"
	"time"

	"github.com/wolfsblu/recipe-manager/domain/security"
)

type UserService struct {
	sender NotificationSender
	store  UserStore
}

func (s *UserService) ConfirmUserByToken(ctx context.Context, token string) error {
	registration, err := s.store.GetRegistrationByToken(ctx, token)
	if err != nil {
		return err
	}

	user := registration.User
	user.Confirmed = true

	return s.store.ConfirmRegistration(ctx, user)
}

func (s *UserService) DeletePasswordResetsOlderThan(ctx context.Context, olderThan time.Duration) error {
	before := time.Now().Add(-olderThan)
	return s.store.DeleteRegistrationsBefore(ctx, before)
}

func (s *UserService) DeleteRegistrationsOlderThan(ctx context.Context, olderThan time.Duration) error {
	before := time.Now().Add(-olderThan)
	return s.store.DeletePasswordResetsBefore(ctx, before)
}

func (s *UserService) UpdatePasswordByToken(ctx context.Context, searchToken, hashedPassword string) error {
	return s.store.UpdatePasswordByToken(ctx, searchToken, hashedPassword)
}

func (s *UserService) GetUserById(ctx context.Context, id int64) (User, error) {
	return s.store.GetUserById(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return s.store.GetUserByEmail(ctx, email)
}

func (s *UserService) RegisterUser(ctx context.Context, userDetails UserDetails) error {
	_, err := s.store.GetUserByEmail(ctx, userDetails.Email)
	if err == nil {
		return ErrUserExists
	}

	_, registration, err := s.store.RegisterUser(ctx, userDetails)
	if err != nil {
		return err
	}

	go func() {
		_ = s.sender.SendUserRegistration(registration)
	}()
	return nil
}

func (s *UserService) ResetPasswordByEmail(ctx context.Context, email string) error {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	} else if !user.Confirmed {
		return ErrUnconfirmedUser
	}
	token, err := s.store.GetPasswordResetTokenByUser(ctx, &user)
	if err == nil { // When there already is a reset token we want to do nothing and exit early
		return nil
	}
	token, err = s.store.CreatePasswordResetToken(ctx, &user)
	if err != nil {
		return err
	}
	go func() {
		_ = s.sender.SendPasswordReset(token)
	}()
	return nil
}

func (s *UserService) VerifyPassword(user User, password string) error {
	ok, err := security.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return WrapError(ErrInvalidCredentials, err)
	} else if !ok {
		return ErrInvalidCredentials
	}
	return nil
}
