package sqlite

import (
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/domain/permissions"
)

func (r *GetPermissionsByRoleRow) AsDomainModel() domain.Permission {
	return domain.Permission{
		ID:   r.ID,
		Name: r.Name,
		Slug: permissions.Slug(r.Slug),
	}
}

func (r *GetUserRow) AsDomainModel() domain.User {
	return domain.User{
		ID:        r.ID,
		Confirmed: r.IsConfirmed,
		Credentials: domain.Credentials{
			Email:        r.Email,
			PasswordHash: r.PasswordHash,
		},
		Role: domain.Role{
			ID:   r.ID,
			Name: r.RoleName,
		},
	}
}

func (t *PasswordReset) AsDomainModel() domain.PasswordResetToken {
	return domain.PasswordResetToken{
		Token:     t.Token,
		CreatedAt: t.CreatedAt,
	}
}

func (r *Ingredient) AsDomainModel() domain.Ingredient {
	return domain.Ingredient{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (r *Recipe) AsDomainModel() domain.Recipe {
	return domain.Recipe{
		ID: r.ID,
		RecipeDetails: domain.RecipeDetails{
			Name:        r.Name,
			Description: r.Description,
			CreatedBy: &domain.User{
				ID: r.CreatedBy,
			},
			Servings: r.Servings,
			Minutes:  r.Minutes,
		},
	}
}

func (r *User) AsDomainModel() domain.User {
	return domain.User{
		ID:        r.ID,
		Confirmed: r.IsConfirmed,
		Credentials: domain.Credentials{
			Email:        r.Email,
			PasswordHash: r.PasswordHash,
		},
	}
}

func (r *UserRegistration) AsDomainModel() domain.UserRegistration {
	return domain.UserRegistration{
		Token:     r.Token,
		CreatedAt: r.CreatedAt,
	}
}

func (u *Unit) AsDomainModel() domain.Unit {
	return domain.Unit{
		ID:   u.ID,
		Name: u.Name,
		Code: u.Code,
	}
}
