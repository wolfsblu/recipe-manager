package domain

import (
	"context"
	"time"

	"github.com/wolfsblu/recipe-manager/domain/pagination"
)

type RecipeService struct {
	sender NotificationSender
	store  RecipeStore
}

func (s *RecipeService) Add(ctx context.Context, r Recipe) (Recipe, error) {
	if err := s.validateRecipe(ctx, r); err != nil {
		return Recipe{}, err
	}
	return s.store.CreateRecipe(ctx, r)
}

func (s *RecipeService) GetMealPlan(ctx context.Context, user *User, from time.Time, until time.Time, req pagination.Page) (pagination.Result[MealPlan], error) {
	return s.store.GetMealPlan(ctx, user, from, until, req)
}

func (s *RecipeService) CreateMealPlan(ctx context.Context, user *User, recipeID int64, date time.Time) error {
	// Use a simple pagination request to check existing entries (just fetch first page)
	existingMealPlan, err := s.store.GetMealPlan(ctx, user, date, date, pagination.Page{Cursor: "", Limit: pagination.DefaultLimit})
	if err != nil {
		return err
	}

	var sortOrder int64 = 0
	if len(existingMealPlan.Data) > 0 && len(existingMealPlan.Data[0].Recipes) > 0 {
		sortOrder = int64(len(existingMealPlan.Data[0].Recipes))
	}

	return s.store.CreateMealPlan(ctx, MealPlanEntry{
		UserID:    user.ID,
		RecipeID:  recipeID,
		Date:      date,
		SortOrder: sortOrder,
	})
}

func (s *RecipeService) DeleteMealPlan(ctx context.Context, user *User, recipeID int64, date time.Time) error {
	return s.store.DeleteMealPlan(ctx, user.ID, recipeID, date)
}

func (s *RecipeService) Delete(ctx context.Context, user *User, id int64) error {
	if err := s.validateRecipeOwnership(ctx, user, id); err != nil {
		return err
	}
	return s.store.DeleteRecipe(ctx, id)
}

func (s *RecipeService) GetByUser(ctx context.Context, user *User, req pagination.Page) (pagination.Result[Recipe], error) {
	return s.store.GetRecipesByUser(ctx, user, req)
}

func (s *RecipeService) GetById(ctx context.Context, user *User, id int64) (Recipe, error) {
	return s.store.GetRecipeById(ctx, user, id)
}

func (s *RecipeService) GetIngredients(ctx context.Context, req pagination.Page) (pagination.Result[Ingredient], error) {
	return s.store.GetIngredients(ctx, req)
}

func (s *RecipeService) GetUnits(ctx context.Context, req pagination.Page) (pagination.Result[Unit], error) {
	return s.store.GetUnits(ctx, req)
}

func (s *RecipeService) GetTags(ctx context.Context, req pagination.Page) (pagination.Result[Tag], error) {
	return s.store.GetTags(ctx, req)
}

func (s *RecipeService) UpdateRecipe(ctx context.Context, recipe Recipe) (Recipe, error) {
	if err := s.validateRecipeOwnership(ctx, recipe.CreatedBy, recipe.ID); err != nil {
		return Recipe{}, err
	}

	if err := s.validateRecipe(ctx, recipe); err != nil {
		return Recipe{}, err
	}

	return s.store.UpdateRecipe(ctx, recipe)
}

func (s *RecipeService) AddIngredient(ctx context.Context, ingredient Ingredient) (Ingredient, error) {
	if err := s.validateIngredient(ingredient); err != nil {
		return Ingredient{}, err
	}
	return s.store.CreateIngredient(ctx, ingredient)
}

func (s *RecipeService) UpdateIngredient(ctx context.Context, ingredient Ingredient) (Ingredient, error) {
	if err := s.validateIngredient(ingredient); err != nil {
		return Ingredient{}, err
	}
	return s.store.UpdateIngredient(ctx, ingredient)
}

func (s *RecipeService) DeleteIngredient(ctx context.Context, id int64) error {
	return s.store.DeleteIngredient(ctx, id)
}

func (s *RecipeService) AddUnit(ctx context.Context, unit Unit) (Unit, error) {
	if err := s.validateUnit(unit); err != nil {
		return Unit{}, err
	}
	return s.store.CreateUnit(ctx, unit)
}

func (s *RecipeService) UpdateUnit(ctx context.Context, unit Unit) (Unit, error) {
	if err := s.validateUnit(unit); err != nil {
		return Unit{}, err
	}
	if err := s.store.UpdateUnit(ctx, unit); err != nil {
		return Unit{}, err
	}
	return unit, nil
}

func (s *RecipeService) DeleteUnit(ctx context.Context, id int64) error {
	return s.store.DeleteUnit(ctx, id)
}
