package domain

import (
	"context"
)

func (s *RecipeService) validateRecipe(ctx context.Context, r Recipe) error {
	return nil
}

func (s *RecipeService) validateRecipeOwnership(ctx context.Context, user *User, recipeID int64) error {
	recipe, err := s.store.GetRecipeById(ctx, user, recipeID)
	if err != nil {
		return err
	}

	if recipe.CreatedBy.ID != user.ID {
		return ErrAuthorization
	}

	return nil
}

func (s *RecipeService) validateVote(vote int64) error {
	if vote != 1 && vote != -1 {
		return ErrInvalidVote
	}
	return nil
}

func (s *RecipeService) validateIngredient(ingredient Ingredient) error {
	if ingredient.Name == "" {
		return ErrInvalidIngredient
	}
	return nil
}

func (s *RecipeService) validateUnit(unit Unit) error {
	if unit.Name == "" {
		return ErrInvalidUnit
	}
	return nil
}
