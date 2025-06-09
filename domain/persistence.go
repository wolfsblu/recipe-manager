package domain

import (
	"context"
	"time"
)

type RecipeStore interface {
	Begin(ctx context.Context) error
	BrowseRecipes(ctx context.Context) ([]Recipe, error)
	Commit() error
	CreateRecipe(ctx context.Context, recipe RecipeDetails) (Recipe, error)
	CreatePasswordResetToken(ctx context.Context, user *User) (PasswordResetToken, error)
	CreateUser(ctx context.Context, credentials Credentials) (User, error)
	CreateUserRegistration(ctx context.Context, user *User) (UserRegistration, error)
	DeletePasswordResetsBefore(ctx context.Context, before time.Time) error
	DeleteRecipe(ctx context.Context, id int64) error
	DeleteRegistrationByUser(ctx context.Context, user *User) error
	DeleteRegistrationsBefore(ctx context.Context, before time.Time) error
	GetMealPlan(ctx context.Context, user *User, from time.Time, until time.Time) ([]MealPlan, error)
	GetPasswordResetTokenByUser(ctx context.Context, user *User) (PasswordResetToken, error)
	GetRecipeById(ctx context.Context, id int64) (Recipe, error)
	GetRecipesByUser(ctx context.Context, user *User) ([]Recipe, error)
	GetRegistrationByToken(ctx context.Context, token string) (UserRegistration, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	Rollback()
	UpdatePasswordByToken(ctx context.Context, token, hashedPassword string) error
	UpdateUser(ctx context.Context, user *User) error
}
