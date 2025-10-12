package mapper

import (
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (m *DBMapper) FromRecipe(recipe domain.Recipe) database.CreateRecipeParams {
	return database.CreateRecipeParams{
		Name:        recipe.Name,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Description: recipe.Description,
		CreatedBy:   recipe.CreatedBy.ID,
	}
}

func (m *DBMapper) FromRecipeForUpdate(recipe domain.Recipe) database.UpdateRecipeParams {
	return database.UpdateRecipeParams{
		Name:        recipe.Name,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Description: recipe.Description,
		ID:          recipe.ID,
	}
}

func (m *DBMapper) FromRecipeStep(recipeID int64, step domain.RecipeStep, sortOrder int64) database.CreateRecipeStepParams {
	return database.CreateRecipeStepParams{
		RecipeID:     recipeID,
		Instructions: step.Instructions,
		SortOrder:    sortOrder,
	}
}

func (m *DBMapper) FromStepIngredient(stepID int64, ingredient domain.StepIngredient, sortOrder int64) database.CreateStepIngredientParams {
	return database.CreateStepIngredientParams{
		StepID:       stepID,
		IngredientID: ingredient.Ingredient.ID,
		UnitID:       ingredient.Unit.ID,
		Amount:       ingredient.Amount,
		SortOrder:    sortOrder,
	}
}

func (m *DBMapper) FromRecipeImage(recipeID int64, image domain.RecipeImage) database.CreateRecipeImagesParams {
	return database.CreateRecipeImagesParams{
		RecipeID: recipeID,
		Path:     image.URL.String(),
	}
}

func (m *DBMapper) FromRecipeTag(recipeID int64, tag domain.Tag) database.CreateRecipeTagParams {
	return database.CreateRecipeTagParams{
		RecipeID: recipeID,
		TagID:    tag.ID,
	}
}

func (m *DBMapper) FromUserVotesParams(recipeIds []int64, userID int64) database.GetUserVotesForRecipesParams {
	return database.GetUserVotesForRecipesParams{
		RecipeIds: recipeIds,
		UserID:    userID,
	}
}

func (m *DBMapper) FromMealPlanEntry(entry domain.MealPlanEntry) database.CreateMealPlanParams {
	return database.CreateMealPlanParams{
		Date:      entry.Date.Format(time.DateOnly),
		UserID:    entry.UserID,
		RecipeID:  entry.RecipeID,
		SortOrder: entry.SortOrder,
	}
}

func (m *DBMapper) FromIngredientNutrient(ingredientID int64, nutrient domain.IngredientNutrient) database.AddIngredientNutrientParams {
	return database.AddIngredientNutrientParams{
		IngredientID: ingredientID,
		NutrientID:   nutrient.Nutrient.ID,
		Amount:       nutrient.Amount,
	}
}
