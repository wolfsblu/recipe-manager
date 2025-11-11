package mapper

import (
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/domain/permissions"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/gen/model"
)

func (m *DBMapper) ToPermission(r model.Permission) domain.Permission {
	return domain.Permission{
		ID:   r.ID,
		Name: r.Name,
		Slug: permissions.Slug(r.Slug),
	}
}

func (m *DBMapper) ToPasswordResetToken(t model.PasswordReset) domain.PasswordResetToken {
	return domain.PasswordResetToken{
		Token:     t.Token,
		CreatedAt: t.CreatedAt,
	}
}

func (m *DBMapper) ToUser(r model.User) domain.User {
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

func (m *DBMapper) ToUserRegistration(r model.UserRegistration) domain.UserRegistration {
	return domain.UserRegistration{
		Token:     r.Token,
		CreatedAt: r.CreatedAt,
	}
}
