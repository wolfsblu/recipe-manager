package domain

import (
	"context"
)

func (s *RecipeService) validateRecipe(ctx context.Context, r Recipe) error {
	return nil
}

func (s *RecipeService) validateRecipeOwnership(ctx context.Context, userID int64, recipeID int64) error {
	recipe, err := s.store.GetRecipeById(ctx, recipeID)
	if err != nil {
		return err
	}

	if recipe.CreatedBy.ID != userID {
		return ErrAuthorization
	}

	return nil
}
