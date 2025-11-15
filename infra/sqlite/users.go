package sqlite

import (
	"context"
	"log"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/domain/roles"
	"github.com/wolfsblu/recipe-manager/domain/security"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/gen/model"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/queries"
)

func (s *Store) CreatePasswordResetToken(ctx context.Context, user *domain.User) (token domain.PasswordResetToken, _ error) {
	generatedToken := security.GenerateToken(security.DefaultTokenLength)

	var result model.PasswordReset
	err := queries.InsertPasswordReset(user.ID, generatedToken).QueryContext(ctx, s.DB(), &result)
	if err != nil {
		return token, domain.WrapError(domain.ErrCreatingPasswordResetToken, err)
	}

	token = s.mapper.ToPasswordResetToken(result)
	token.User = user
	return token, nil
}

func (s *Store) DeletePasswordResetsBefore(ctx context.Context, before time.Time) error {
	_, err := queries.DeletePasswordResetsBefore(before).ExecContext(ctx, s.DB())
	return err
}

func (s *Store) DeleteRegistrationsBefore(ctx context.Context, before time.Time) error {
	_, err := queries.DeleteRegistrationsBefore(before).ExecContext(ctx, s.DB())
	return err
}

func (s *Store) GetPasswordResetTokenByUser(ctx context.Context, user *domain.User) (token domain.PasswordResetToken, _ error) {
	var result model.PasswordReset
	err := queries.SelectPasswordResetByUser(user.ID).QueryContext(ctx, s.DB(), &result)
	if err != nil {
		return token, domain.WrapError(domain.ErrPasswordResetTokenNotFound, err)
	}
	token = s.mapper.ToPasswordResetToken(result)
	token.User = user
	return token, nil
}

func (s *Store) GetRegistrationByToken(ctx context.Context, token string) (registration domain.UserRegistration, _ error) {
	type ResultRow struct {
		UserRegistration model.UserRegistration
		User             model.User
	}
	var result ResultRow
	err := queries.SelectRegistrationByToken(token).QueryContext(ctx, s.DB(), &result)
	if err != nil {
		return registration, domain.WrapError(domain.ErrRegistrationNotFound, err)
	}
	user := s.mapper.ToUser(result.User)
	registration = s.mapper.ToUserRegistration(result.UserRegistration)
	registration.User = &user
	return registration, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (user domain.User, _ error) {
	var result model.User
	log.Println(email)
	stmt := queries.SelectUserByEmail(email)
	log.Println(stmt.DebugSql())
	err := stmt.QueryContext(ctx, s.DB(), &result)

	if err != nil {
		return user, domain.WrapError(domain.ErrUserNotFound, err)
	}
	return s.mapper.ToUser(result), nil
}

func (s *Store) GetUserById(ctx context.Context, id int64) (user domain.User, _ error) {
	type ResultRow struct {
		model.User
		Role struct {
			Name string
		}
	}
	var result ResultRow
	err := queries.SelectUserByID(id).QueryContext(ctx, s.DB(), &result)
	if err != nil {
		return user, domain.WrapError(domain.ErrUserNotFound, err)
	}
	user = s.mapper.ToUser(result.User)
	user.Role = domain.Role{
		ID:   result.User.RoleID,
		Name: result.Role.Name,
	}

	// Get permissions
	var permissions []model.Permission
	err = queries.SelectPermissionsByRole(result.User.RoleID).QueryContext(ctx, s.DB(), &permissions)
	if err != nil {
		return user, err
	}

	domainPermissions := make([]domain.Permission, len(permissions))
	for i, perm := range permissions {
		domainPermissions[i] = s.mapper.ToPermission(perm)
	}
	user.Role.Permissions = domainPermissions

	return user, nil
}

func (s *Store) UpdatePasswordByToken(ctx context.Context, searchToken, hashedPassword string) error {
	return s.WithTransaction(ctx, func(tx *TxStore) error {
		type ResultRow struct {
			PasswordReset model.PasswordReset
			User          model.User
		}
		var result ResultRow
		err := queries.SelectPasswordResetByToken(searchToken).QueryContext(ctx, tx.DB(), &result)
		if err != nil {
			return domain.WrapError(domain.ErrPasswordResetTokenNotFound, err)
		}

		_, err = queries.UpdatePassword(result.User.ID, hashedPassword).ExecContext(ctx, tx.DB())
		if err != nil {
			return domain.WrapError(domain.ErrUpdatingPassword, err)
		}

		_, err = queries.DeletePasswordReset(result.User.ID).ExecContext(ctx, tx.DB())
		if err != nil {
			return domain.WrapError(domain.ErrDeletingPasswordResetToken, err)
		}

		return nil
	})
}

func (s *Store) ConfirmRegistration(ctx context.Context, user *domain.User) error {
	return s.WithTransaction(ctx, func(tx *TxStore) error {
		_, err := queries.UpdateUser(user.ID, user.Email, user.Confirmed).ExecContext(ctx, tx.DB())
		if err != nil {
			return domain.WrapError(domain.ErrUpdatingUser, err)
		}

		_, err = queries.DeleteRegistration(user.ID).ExecContext(ctx, tx.DB())
		if err != nil {
			return domain.WrapError(domain.ErrDeletingRegistration, err)
		}

		return nil
	})
}

func (s *Store) RegisterUser(ctx context.Context, userDetails domain.UserDetails) (domain.User, domain.UserRegistration, error) {
	var user domain.User
	var registration domain.UserRegistration

	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		var dbUser model.User
		err := queries.InsertUser(userDetails.Email, userDetails.PasswordHash, int64(roles.User), userDetails.Locale).
			QueryContext(ctx, tx.DB(), &dbUser)
		if err != nil {
			return err
		}
		user = s.mapper.ToUser(dbUser)

		generatedToken := security.GenerateToken(security.DefaultTokenLength)
		var dbRegistration model.UserRegistration
		err = queries.InsertUserRegistration(user.ID, generatedToken).QueryContext(ctx, tx.DB(), &dbRegistration)
		if err != nil {
			return err
		}
		registration = s.mapper.ToUserRegistration(dbRegistration)

		return nil
	})

	return user, registration, err
}
