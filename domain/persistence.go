package domain

import (
	"context"
	"time"
)

type Transactioner interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback()
}

type RecipeStore interface {
	Transactioner
	BrowseRecipes(ctx context.Context) ([]Recipe, error)
	CreateRecipe(ctx context.Context, recipe Recipe) (Recipe, error)
	DeleteRecipe(ctx context.Context, id int64) error
	GetMealPlan(ctx context.Context, user *User, from time.Time, until time.Time) ([]MealPlan, error)
	GetIngredients(ctx context.Context) ([]Ingredient, error)
	GetUnits(ctx context.Context) ([]Unit, error)
	GetTags(ctx context.Context) ([]Tag, error)
	GetRecipeById(ctx context.Context, id int64) (Recipe, error)
	GetRecipesByUser(ctx context.Context, user *User) ([]Recipe, error)
	UpdateRecipe(ctx context.Context, recipe Recipe) (Recipe, error)
}

type UserStore interface {
	Transactioner
	CreatePasswordResetToken(ctx context.Context, user *User) (PasswordResetToken, error)
	CreateUser(ctx context.Context, credentials Credentials) (User, error)
	CreateUserRegistration(ctx context.Context, user *User) (UserRegistration, error)
	DeletePasswordResetsBefore(ctx context.Context, before time.Time) error
	DeleteRegistrationByUser(ctx context.Context, user *User) error
	DeleteRegistrationsBefore(ctx context.Context, before time.Time) error
	GetPasswordResetTokenByUser(ctx context.Context, user *User) (PasswordResetToken, error)
	GetRegistrationByToken(ctx context.Context, token string) (UserRegistration, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	UpdatePasswordByToken(ctx context.Context, token, hashedPassword string) error
	UpdateUser(ctx context.Context, user *User) error
}
