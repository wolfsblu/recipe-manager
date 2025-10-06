package sqlite

import (
	"context"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/domain/roles"
	"github.com/wolfsblu/recipe-manager/domain/security"
)

func (s *Store) CreatePasswordResetToken(ctx context.Context, user *domain.User) (token domain.PasswordResetToken, _ error) {
	generatedToken := security.GenerateToken(security.DefaultTokenLength)
	result, err := s.query().CreatePasswordResetToken(ctx, s.mapper.FromUserForPasswordResetToken(user, generatedToken))
	if err != nil {
		return token, domain.WrapError(domain.ErrCreatingPasswordResetToken, err)
	}

	token = s.mapper.ToPasswordResetToken(result)
	token.User = user
	return token, nil
}

func (s *Store) DeletePasswordResetsBefore(ctx context.Context, before time.Time) error {
	return s.query().DeletePasswordResetsBefore(ctx, before)
}

func (s *Store) DeleteRegistrationsBefore(ctx context.Context, before time.Time) error {
	return s.query().DeleteRegistrationsBefore(ctx, before)
}

func (s *Store) GetPasswordResetTokenByUser(ctx context.Context, user *domain.User) (token domain.PasswordResetToken, _ error) {
	result, err := s.query().GetPasswordResetTokenByUser(ctx, user.ID)
	if err != nil {
		return token, domain.WrapError(domain.ErrPasswordResetTokenNotFound, err)
	}
	token = s.mapper.ToPasswordResetToken(result)
	token.User = user
	return token, nil
}

func (s *Store) GetRegistrationByToken(ctx context.Context, token string) (registration domain.UserRegistration, _ error) {
	result, err := s.query().GetUserRegistration(ctx, token)
	if err != nil {
		return registration, domain.WrapError(domain.ErrRegistrationNotFound, err)
	}
	user := s.mapper.ToUser(result.User)
	registration = s.mapper.ToUserRegistration(result.UserRegistration)
	registration.User = &user
	return registration, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (user domain.User, _ error) {
	result, err := s.query().GetUserByEmail(ctx, email)
	if err != nil {
		return user, domain.WrapError(domain.ErrUserNotFound, err)
	}
	return s.mapper.ToUser(result), nil
}

func (s *Store) GetUserById(ctx context.Context, id int64) (user domain.User, _ error) {
	result, err := s.query().GetUser(ctx, id)
	if err != nil {
		return user, domain.WrapError(domain.ErrUserNotFound, err)
	}
	user = s.mapper.ToUserFromGetUserRow(result)

	permissionResult, err := s.query().GetPermissionsByRole(ctx, result.RoleID)
	permissions := make([]domain.Permission, len(permissionResult))
	for i, permissionRow := range permissionResult {
		permissions[i] = s.mapper.ToPermission(permissionRow)
	}
	user.Role.Permissions = permissions

	return user, nil
}

func (s *Store) UpdatePasswordByToken(ctx context.Context, searchToken, hashedPassword string) error {
	return s.WithTransaction(ctx, func(tx *TxStore) error {
		token, err := tx.query().GetPasswordResetToken(ctx, searchToken)
		if err != nil {
			return domain.WrapError(domain.ErrPasswordResetTokenNotFound, err)
		}
		if err = tx.query().UpdatePasswordByUserId(ctx, s.mapper.FromUserForPasswordUpdate(hashedPassword, token.User.ID)); err != nil {
			return domain.WrapError(domain.ErrUpdatingPassword, err)
		}
		if err = tx.query().DeletePasswordResetTokenByUserId(ctx, token.User.ID); err != nil {
			return domain.WrapError(domain.ErrDeletingPasswordResetToken, err)
		}
		return nil
	})
}

func (s *Store) ConfirmRegistration(ctx context.Context, user *domain.User) error {
	return s.WithTransaction(ctx, func(tx *TxStore) error {
		if err := tx.query().UpdateUser(ctx, s.mapper.FromUserForUpdate(user)); err != nil {
			return domain.WrapError(domain.ErrUpdatingUser, err)
		}
		if err := tx.query().DeleteRegistrationByUserId(ctx, user.ID); err != nil {
			return domain.WrapError(domain.ErrDeletingRegistration, err)
		}
		return nil
	})
}

func (s *Store) RegisterUser(ctx context.Context, userDetails domain.UserDetails) (domain.User, domain.UserRegistration, error) {
	var user domain.User
	var registration domain.UserRegistration

	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		var err error
		userParams := s.mapper.FromUserDetails(userDetails, int64(roles.User))
		dbUser, err := tx.query().CreateUser(ctx, userParams)
		user = s.mapper.ToUser(dbUser)
		if err != nil {
			return err
		}

		generatedToken := security.GenerateToken(security.DefaultTokenLength)
		registrationParams := s.mapper.FromUserForRegistration(&user, generatedToken)
		dbRegistration, err := tx.query().CreateUserRegistration(ctx, registrationParams)
		registration = s.mapper.ToUserRegistration(dbRegistration)
		if err != nil {
			return err
		}

		return nil
	})

	return user, registration, err
}
