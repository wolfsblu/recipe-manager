package domain

import (
	"context"
	"github.com/wolfsblu/go-chef/domain/security"
	"time"
)

type Credentials struct {
	Email        string
	PasswordHash string
}

type User struct {
	ID        int64
	Confirmed bool
	Credentials
}

type PasswordResetToken struct {
	User      *User
	Token     string
	CreatedAt time.Time
}

type UserRegistration struct {
	User      *User
	Token     string
	CreatedAt time.Time
}

func (s *RecipeService) ConfirmUserByToken(ctx context.Context, token string) error {
	registration, err := s.store.GetRegistrationByToken(ctx, token)
	if err != nil {
		return err
	}

	user := registration.User
	user.Confirmed = true

	if err := s.store.Begin(ctx); err != nil {
		return err
	}
	defer s.store.Rollback()
	if err = s.store.UpdateUser(ctx, user); err != nil {
		return err
	}
	if err = s.store.DeleteRegistrationByUser(ctx, user); err != nil {
		return err
	}
	return s.store.Commit()
}

func (s *RecipeService) UpdatePasswordByToken(ctx context.Context, searchToken, hashedPassword string) error {
	return s.store.UpdatePasswordByToken(ctx, searchToken, hashedPassword)
}

func (s *RecipeService) GetUserById(ctx context.Context, id int64) (User, error) {
	return s.store.GetUserById(ctx, id)
}

func (s *RecipeService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return s.store.GetUserByEmail(ctx, email)
}

func (s *RecipeService) RegisterUser(ctx context.Context, credentials Credentials) error {
	user, err := s.store.GetUserByEmail(ctx, credentials.Email)
	if err == nil {
		return ErrUserExists
	}
	if err := s.store.Begin(ctx); err != nil {
		return err
	}
	defer s.store.Rollback()
	user, err = s.store.CreateUser(ctx, credentials)
	if err != nil {
		return err
	}
	registration, err := s.store.CreateUserRegistration(ctx, &user)
	if err != nil {
		return err
	}
	if err = s.store.Commit(); err != nil {
		return err
	}
	go func() {
		_ = s.sender.SendUserRegistration(registration)
	}()
	return nil
}

func (s *RecipeService) ResetPasswordByEmail(ctx context.Context, email string) error {
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return err
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

func (s *RecipeService) VerifyPassword(user User, password string) error {
	ok, err := security.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return WrapError(ErrInvalidCredentials, err)
	} else if !ok {
		return ErrInvalidCredentials
	}
	return nil
}
