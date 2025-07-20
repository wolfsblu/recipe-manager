package domain

import (
	"context"
	"time"
)

type RecipeService struct {
	sender NotificationSender
	store  RecipeStore
}

type RecipeDetails struct {
	Name        string
	Description string
	CreatedBy   *User
	Servings    int64
	Minutes     int64
}

type Recipe struct {
	ID     int64
	Tags   []string
	Images []RecipeImage
	RecipeDetails
}

type RecipeImage struct {
	ID        int64
	Path      string
	SortOrder int64
}

type MealPlan struct {
	Date    time.Time
	Recipes []Recipe
}

func (s *RecipeService) Add(ctx context.Context, r RecipeDetails) (Recipe, error) {
	return s.store.CreateRecipe(ctx, r)
}

func (s *RecipeService) Browse(ctx context.Context) ([]Recipe, error) {
	return s.store.BrowseRecipes(ctx)
}

func (s *RecipeService) GetMealPlan(ctx context.Context, user *User, from time.Time, until time.Time) ([]MealPlan, error) {
	return s.store.GetMealPlan(ctx, user, from, until)
}

func (s *RecipeService) Delete(ctx context.Context, id int64) error {
	return s.store.DeleteRecipe(ctx, id)
}

func (s *RecipeService) GetByUser(ctx context.Context, user *User) ([]Recipe, error) {
	return s.store.GetRecipesByUser(ctx, user)
}

func (s *RecipeService) GetById(ctx context.Context, id int64) (Recipe, error) {
	return s.store.GetRecipeById(ctx, id)
}
