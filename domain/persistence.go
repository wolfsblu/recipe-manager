package domain

import (
	"context"
	"time"
)

type RecipeStore interface {
	BrowseRecipes(ctx context.Context) ([]Recipe, error)
	CreateRecipe(ctx context.Context, recipe Recipe) (Recipe, error)
	DeleteRecipe(ctx context.Context, id int64) error
	GetMealPlan(ctx context.Context, user *User, from time.Time, until time.Time) ([]MealPlan, error)
	CreateMealPlan(ctx context.Context, entry MealPlanEntry) error
	DeleteMealPlan(ctx context.Context, userID int64, recipeID int64, date time.Time) error
	GetIngredients(ctx context.Context) ([]Ingredient, error)
	GetUnits(ctx context.Context) ([]Unit, error)
	GetTags(ctx context.Context) ([]Tag, error)
	GetRecipeById(ctx context.Context, user *User, id int64) (Recipe, error)
	GetRecipesByUser(ctx context.Context, user *User) ([]Recipe, error)
	UpdateRecipe(ctx context.Context, recipe Recipe) (Recipe, error)
	AddVote(ctx context.Context, recipeID int64, userID int64, vote int64) error
	RemoveVote(ctx context.Context, recipeID int64, userID int64) error
	GetRecipeVotes(ctx context.Context, recipeID int64) (int64, error)
	GetUserVote(ctx context.Context, recipeID int64, userID int64) (int64, error)
	CreateIngredient(ctx context.Context, ingredient Ingredient) (Ingredient, error)
	UpdateIngredient(ctx context.Context, ingredient Ingredient) error
	DeleteIngredient(ctx context.Context, id int64) error
	CreateUnit(ctx context.Context, unit Unit) (Unit, error)
	UpdateUnit(ctx context.Context, unit Unit) error
	DeleteUnit(ctx context.Context, id int64) error
}

type UserStore interface {
	ConfirmRegistration(ctx context.Context, user *User) error
	CreatePasswordResetToken(ctx context.Context, user *User) (PasswordResetToken, error)
	DeletePasswordResetsBefore(ctx context.Context, before time.Time) error
	DeleteRegistrationsBefore(ctx context.Context, before time.Time) error
	GetPasswordResetTokenByUser(ctx context.Context, user *User) (PasswordResetToken, error)
	GetRegistrationByToken(ctx context.Context, token string) (UserRegistration, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	RegisterUser(ctx context.Context, userDetails UserDetails) (User, UserRegistration, error)
	UpdatePasswordByToken(ctx context.Context, token, hashedPassword string) error
}

type ShoppingStore interface {
	GetShoppingListsByUser(ctx context.Context, userID int64) ([]ShoppingList, error)
	GetShoppingListByID(ctx context.Context, listID int64) (ShoppingList, error)
	CreateShoppingList(ctx context.Context, userID int64, name string) (ShoppingList, error)
	UpdateShoppingList(ctx context.Context, listID int64, name string) (ShoppingList, error)
	DeleteShoppingList(ctx context.Context, listID int64) error
	CreateShoppingListItem(ctx context.Context, listID int64, item ShoppingListItem) (ShoppingListItem, error)
	UpdateShoppingListItem(ctx context.Context, itemID int64, item ShoppingListItem) (ShoppingListItem, error)
	DeleteShoppingListItem(ctx context.Context, itemID int64) error
}
