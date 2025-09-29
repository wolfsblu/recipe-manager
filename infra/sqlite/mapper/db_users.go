package mapper

import (
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/domain/permissions"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (m *DBMapper) ToPermission(r database.GetPermissionsByRoleRow) domain.Permission {
	return domain.Permission{
		ID:   r.ID,
		Name: r.Name,
		Slug: permissions.Slug(r.Slug),
	}
}

func (m *DBMapper) ToUserFromGetUserRow(r database.GetUserRow) domain.User {
	return domain.User{
		ID:        r.ID,
		Confirmed: r.IsConfirmed,
		UserDetails: domain.UserDetails{
			Email:        r.Email,
			PasswordHash: r.PasswordHash,
			Locale:       r.Locale,
		},
		Role: domain.Role{
			ID:   r.ID,
			Name: r.RoleName,
		},
	}
}

func (m *DBMapper) ToPasswordResetToken(t database.PasswordReset) domain.PasswordResetToken {
	return domain.PasswordResetToken{
		Token:     t.Token,
		CreatedAt: t.CreatedAt,
	}
}

func (m *DBMapper) ToUser(r database.User) domain.User {
	return domain.User{
		ID:        r.ID,
		Confirmed: r.IsConfirmed,
		UserDetails: domain.UserDetails{
			Email:        r.Email,
			PasswordHash: r.PasswordHash,
			Locale:       r.Locale,
		},
	}
}

func (m *DBMapper) ToUserRegistration(r database.UserRegistration) domain.UserRegistration {
	return domain.UserRegistration{
		Token:     r.Token,
		CreatedAt: r.CreatedAt,
	}
}
