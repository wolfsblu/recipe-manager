package queries

import (
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectRecipeByID returns a query to fetch a single recipe by ID
func SelectRecipeByID(id int64) SelectStatement {
	return SELECT(
		Recipe.ID,
		Recipe.Name,
		Recipe.Servings,
		Recipe.Minutes,
		Recipe.Description,
		Recipe.CreatedBy,
		Recipe.CreatedAt,
	).FROM(Recipe).
		WHERE(Recipe.ID.EQ(Int(id)))
}

// SelectRecipesByUserPaginated returns a query for paginated user recipes with sorting
func SelectRecipesByUserPaginated(userID int64, cursor RecipeCursor, limit int64, sortField, sortOrder string) SelectStatement {
	var comparison BoolExpression
	var orderBy []OrderByClause

	switch sortField {
	case "name":
		if sortOrder == "asc" {
			comparison = ROW(Recipe.Name, Recipe.ID).GT(ROW(String(cursor.LastName), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipe.Name.ASC(), Recipe.ID.ASC()}
		} else {
			comparison = ROW(Recipe.Name, Recipe.ID).LT(ROW(String(cursor.LastName), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipe.Name.DESC(), Recipe.ID.DESC()}
		}
	case "servings":
		if sortOrder == "asc" {
			comparison = ROW(Recipe.Servings, Recipe.ID).GT(ROW(Int(cursor.LastServings), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipe.Servings.ASC(), Recipe.ID.ASC()}
		} else {
			comparison = ROW(Recipe.Servings, Recipe.ID).LT(ROW(Int(cursor.LastServings), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipe.Servings.DESC(), Recipe.ID.DESC()}
		}
	default: // created_at
		if sortOrder == "asc" {
			comparison = ROW(Recipe.CreatedAt, Recipe.ID).GT(ROW(String(cursor.LastCreatedAt.Format(time.RFC3339)), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipe.CreatedAt.ASC(), Recipe.ID.ASC()}
		} else {
			comparison = ROW(Recipe.CreatedAt, Recipe.ID).LT(ROW(String(cursor.LastCreatedAt.Format(time.RFC3339)), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipe.CreatedAt.DESC(), Recipe.ID.DESC()}
		}
	}

	return SELECT(
		Recipe.ID,
		Recipe.Name,
		Recipe.Servings,
		Recipe.Minutes,
		Recipe.Description,
		Recipe.CreatedBy,
		Recipe.CreatedAt,
	).FROM(Recipe).
		WHERE(Recipe.CreatedBy.EQ(Int(userID)).AND(comparison)).
		ORDER_BY(orderBy...).
		LIMIT(limit)
}

// SelectTagsForRecipes returns a query to fetch all tags for given recipe IDs
func SelectTagsForRecipes(recipeIDs []int64) SelectStatement {
	return SELECT(RecipeTag.RecipeID, Tag.AllColumns).FROM(
		Tag.INNER_JOIN(RecipeTag, Tag.ID.EQ(RecipeTag.TagID)),
	).WHERE(RecipeTag.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(Tag.Name.ASC())
}

// SelectImagesForRecipes returns a query to fetch all images for given recipe IDs
func SelectImagesForRecipes(recipeIDs []int64) SelectStatement {
	return SELECT(RecipeImage.AllColumns).FROM(RecipeImage).
		WHERE(RecipeImage.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(RecipeImage.SortOrder.ASC())
}

// InsertRecipe returns an insert statement for a new recipe
func InsertRecipe(name, description string, servings, minutes, createdBy int64) InsertStatement {
	return Recipe.INSERT(Recipe.Name, Recipe.Servings, Recipe.Minutes, Recipe.Description, Recipe.CreatedBy).
		VALUES(name, servings, minutes, description, createdBy).
		RETURNING(Recipe.AllColumns)
}

// DeleteRecipe returns a delete statement for a recipe
func DeleteRecipe(id int64) DeleteStatement {
	return Recipe.DELETE().WHERE(Recipe.ID.EQ(Int(id)))
}

// SelectStepsForRecipes returns a query to fetch all steps for given recipe IDs
func SelectStepsForRecipes(recipeIDs []int64) SelectStatement {
	return SELECT(RecipeStep.AllColumns).FROM(RecipeStep).
		WHERE(RecipeStep.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(RecipeStep.SortOrder.ASC())
}

// SelectIngredientsForRecipes returns a query to fetch all ingredients for given recipe IDs
func SelectIngredientsForRecipes(recipeIDs []int64) SelectStatement {
	return SELECT(
		RecipeIngredient.AllColumns,
		Ingredient.AllColumns,
		Unit.AllColumns,
		RecipeStep.ID.AS("step_id"),
	).FROM(RecipeIngredient.
		INNER_JOIN(RecipeStep, RecipeIngredient.StepID.EQ(RecipeStep.ID)).
		INNER_JOIN(Ingredient, RecipeIngredient.IngredientID.EQ(Ingredient.ID)).
		INNER_JOIN(Unit, RecipeIngredient.UnitID.EQ(Unit.ID)),
	).WHERE(RecipeStep.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(RecipeIngredient.SortOrder.ASC())
}

// SelectNutrientsForRecipes returns a query to fetch nutrients for recipe ingredients
func SelectNutrientsForRecipes(recipeIDs []int64) SelectStatement {
	return SELECT(
		IngredientNutrient.AllColumns,
		Nutrient.AllColumns,
	).FROM(
		IngredientNutrient.
			INNER_JOIN(Nutrient, IngredientNutrient.NutrientID.EQ(Nutrient.ID)).
			INNER_JOIN(RecipeIngredient, IngredientNutrient.IngredientID.EQ(RecipeIngredient.IngredientID)).
			INNER_JOIN(RecipeStep, RecipeIngredient.StepID.EQ(RecipeStep.ID)),
	).WHERE(RecipeStep.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(IngredientNutrient.IngredientID.ASC(), Nutrient.Name.ASC())
}

// SelectTags returns a paginated query for tags
func SelectTags(lastID int64, lastName string, search string, limit int64) SelectStatement {
	var whereClauses []BoolExpression

	// Add pagination cursor
	whereClauses = append(whereClauses, ROW(Tag.Name, Tag.ID).GT(ROW(String(lastName), Int(lastID))))

	// Add search filter if provided
	if search != "" {
		whereClauses = append(whereClauses, Tag.Name.LIKE(String("%"+search+"%")))
	}

	return SELECT(
		Tag.ID,
		Tag.Name,
	).FROM(Tag).
		WHERE(AND(whereClauses...)).
		ORDER_BY(Tag.Name.ASC(), Tag.ID.ASC()).
		LIMIT(limit)
}

// SelectMealPlan returns a query for meal plan entries with recipes
func SelectMealPlan(userID int64, from, until string, lastDate string, lastID int64, limit int64) SelectStatement {
	return SELECT(
		MealPlan.ID,
		MealPlan.Date,
		MealPlan.UserID,
		MealPlan.RecipeID,
		MealPlan.SortOrder,
		Recipe.ID,
		Recipe.Name,
		Recipe.Servings,
		Recipe.Minutes,
		Recipe.Description,
		Recipe.CreatedBy,
		Recipe.CreatedAt,
	).FROM(
		MealPlan.INNER_JOIN(Recipe, Recipe.ID.EQ(MealPlan.RecipeID)),
	).WHERE(
		MealPlan.UserID.EQ(Int(userID)).
			AND(MealPlan.Date.GT_EQ(String(from))).
			AND(MealPlan.Date.LT_EQ(String(until))).
			AND(ROW(MealPlan.Date, MealPlan.ID).GT(ROW(String(lastDate), Int(lastID)))),
	).ORDER_BY(
		MealPlan.Date.ASC(),
		MealPlan.ID.ASC(),
	).LIMIT(limit)
}

// InsertRecipeStep returns an insert statement for a recipe step
func InsertRecipeStep(recipeID int64, instructions string, sortOrder int64) InsertStatement {
	return RecipeStep.INSERT(RecipeStep.RecipeID, RecipeStep.Instructions, RecipeStep.SortOrder).
		VALUES(recipeID, instructions, sortOrder).
		RETURNING(RecipeStep.AllColumns)
}

// InsertStepIngredient returns an insert statement for a step ingredient
func InsertStepIngredient(stepID, ingredientID, unitID int64, amount float64, sortOrder int64) InsertStatement {
	return RecipeIngredient.INSERT(
		RecipeIngredient.StepID,
		RecipeIngredient.IngredientID,
		RecipeIngredient.UnitID,
		RecipeIngredient.Amount,
		RecipeIngredient.SortOrder,
	).VALUES(stepID, ingredientID, unitID, amount, sortOrder)
}

// InsertRecipeImage returns an insert statement for a recipe image
func InsertRecipeImage(recipeID int64, path string, sortOrder int64) InsertStatement {
	return RecipeImage.INSERT(RecipeImage.RecipeID, RecipeImage.Path, RecipeImage.SortOrder).
		VALUES(recipeID, path, sortOrder)
}

// InsertRecipeTag returns an insert statement for a recipe tag
func InsertRecipeTag(recipeID, tagID int64) InsertStatement {
	return RecipeTag.INSERT(RecipeTag.RecipeID, RecipeTag.TagID).
		VALUES(recipeID, tagID)
}

// InsertMealPlan returns an insert statement for a meal plan entry
func InsertMealPlan(date string, userID, recipeID, sortOrder int64) InsertStatement {
	return MealPlan.INSERT(MealPlan.Date, MealPlan.UserID, MealPlan.RecipeID, MealPlan.SortOrder).
		VALUES(date, userID, recipeID, sortOrder)
}

// UpdateRecipe returns an update statement for a recipe
func UpdateRecipe(id int64, name, description string, servings, minutes int64) UpdateStatement {
	return Recipe.UPDATE(Recipe.Name, Recipe.Servings, Recipe.Minutes, Recipe.Description).
		SET(name, servings, minutes, description).
		WHERE(Recipe.ID.EQ(Int(id)))
}

// DeleteRecipeSteps returns a delete statement for recipe steps
func DeleteRecipeSteps(recipeID int64) DeleteStatement {
	return RecipeStep.DELETE().WHERE(RecipeStep.RecipeID.EQ(Int(recipeID)))
}

// DeleteRecipeImages returns a delete statement for recipe images
func DeleteRecipeImages(recipeID int64) DeleteStatement {
	return RecipeImage.DELETE().WHERE(RecipeImage.RecipeID.EQ(Int(recipeID)))
}

// DeleteRecipeTags returns a delete statement for recipe tags
func DeleteRecipeTags(recipeID int64) DeleteStatement {
	return RecipeTag.DELETE().WHERE(RecipeTag.RecipeID.EQ(Int(recipeID)))
}

// DeleteRecipeIngredients returns a delete statement for recipe ingredients
func DeleteRecipeIngredients(recipeID int64) DeleteStatement {
	return RecipeIngredient.DELETE().
		WHERE(RecipeIngredient.StepID.IN(
			SELECT(RecipeStep.ID).
				FROM(RecipeStep).
				WHERE(RecipeStep.RecipeID.EQ(Int(recipeID))),
		))
}

// DeleteMealPlan returns a delete statement for a meal plan entry
func DeleteMealPlan(userID, recipeID int64, date string) DeleteStatement {
	return MealPlan.DELETE().WHERE(
		MealPlan.UserID.EQ(Int(userID)).
			AND(MealPlan.RecipeID.EQ(Int(recipeID))).
			AND(MealPlan.Date.EQ(String(date))),
	)
}

// Helper types and functions

type RecipeCursor struct {
	LastID        int64
	LastName      string
	LastCreatedAt time.Time
	LastServings  int64
}

// IntSliceToExpressions converts a slice of int64 to go-jet expressions
func IntSliceToExpressions(ids []int64) []Expression {
	exprs := make([]Expression, len(ids))
	for i, id := range ids {
		exprs[i] = Int(id)
	}
	return exprs
}
