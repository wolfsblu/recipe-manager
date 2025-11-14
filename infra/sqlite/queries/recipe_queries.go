package queries

import (
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectRecipeByID returns a query to fetch a single recipe by ID
func SelectRecipeByID(id int64) SelectStatement {
	return SELECT(
		Recipes.ID,
		Recipes.Name,
		Recipes.Servings,
		Recipes.Minutes,
		Recipes.Description,
		Recipes.CreatedBy,
		Recipes.CreatedAt,
	).FROM(Recipes).
		WHERE(Recipes.ID.EQ(Int(id)))
}

// SelectRecipesByUserPaginated returns a query for paginated user recipes with sorting
func SelectRecipesByUserPaginated(userID int64, cursor RecipeCursor, limit int64, sortField, sortOrder string) SelectStatement {
	var comparison BoolExpression
	var orderBy []OrderByClause

	switch sortField {
	case "name":
		if sortOrder == "asc" {
			comparison = ROW(Recipes.Name, Recipes.ID).GT(ROW(String(cursor.LastName), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipes.Name.ASC(), Recipes.ID.ASC()}
		} else {
			comparison = ROW(Recipes.Name, Recipes.ID).LT(ROW(String(cursor.LastName), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipes.Name.DESC(), Recipes.ID.DESC()}
		}
	case "servings":
		if sortOrder == "asc" {
			comparison = ROW(Recipes.Servings, Recipes.ID).GT(ROW(Int(cursor.LastServings), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipes.Servings.ASC(), Recipes.ID.ASC()}
		} else {
			comparison = ROW(Recipes.Servings, Recipes.ID).LT(ROW(Int(cursor.LastServings), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipes.Servings.DESC(), Recipes.ID.DESC()}
		}
	default: // created_at
		if sortOrder == "asc" {
			comparison = ROW(Recipes.CreatedAt, Recipes.ID).GT(ROW(String(cursor.LastCreatedAt.Format(time.RFC3339)), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipes.CreatedAt.ASC(), Recipes.ID.ASC()}
		} else {
			comparison = ROW(Recipes.CreatedAt, Recipes.ID).LT(ROW(String(cursor.LastCreatedAt.Format(time.RFC3339)), Int(cursor.LastID)))
			orderBy = []OrderByClause{Recipes.CreatedAt.DESC(), Recipes.ID.DESC()}
		}
	}

	return SELECT(
		Recipes.ID,
		Recipes.Name,
		Recipes.Servings,
		Recipes.Minutes,
		Recipes.Description,
		Recipes.CreatedBy,
		Recipes.CreatedAt,
	).FROM(Recipes).
		WHERE(Recipes.CreatedBy.EQ(Int(userID)).AND(comparison)).
		ORDER_BY(orderBy...).
		LIMIT(limit)
}

// SelectTagsForRecipes returns a query to fetch all tags for given recipe IDs
func SelectTagsForRecipes(recipeIDs []int64) SelectStatement {
	tagCols := Tags.AllColumns()
	cols := make([]Projection, 0, len(tagCols)+1)
	cols = append(cols, RecipeTags.RecipeID)
	cols = append(cols, tagCols...)

	return SELECT(cols[0], cols[1:]...).FROM(
		Tags.INNER_JOIN(RecipeTags, Tags.ID.EQ(RecipeTags.TagID)),
	).WHERE(RecipeTags.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(Tags.Name.ASC())
}

// SelectImagesForRecipes returns a query to fetch all images for given recipe IDs
func SelectImagesForRecipes(recipeIDs []int64) SelectStatement {
	cols := RecipeImages.AllColumns()
	return SELECT(cols[0], cols[1:]...).FROM(RecipeImages).
		WHERE(RecipeImages.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(RecipeImages.SortOrder.ASC())
}

// InsertRecipe returns an insert statement for a new recipe
func InsertRecipe(name, description string, servings, minutes, createdBy int64) InsertStatement {
	return Recipes.INSERT(Recipes.Name, Recipes.Servings, Recipes.Minutes, Recipes.Description, Recipes.CreatedBy).
		VALUES(name, servings, minutes, description, createdBy).
		RETURNING(Recipes.AllColumns()...)
}

// DeleteRecipe returns a delete statement for a recipe
func DeleteRecipe(id int64) DeleteStatement {
	return Recipes.DELETE().WHERE(Recipes.ID.EQ(Int(id)))
}

// SelectStepsForRecipes returns a query to fetch all steps for given recipe IDs
func SelectStepsForRecipes(recipeIDs []int64) SelectStatement {
	cols := RecipeSteps.AllColumns()
	return SELECT(cols[0], cols[1:]...).FROM(RecipeSteps).
		WHERE(RecipeSteps.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(RecipeSteps.SortOrder.ASC())
}

// SelectIngredientsForRecipes returns a query to fetch all ingredients for given recipe IDs
func SelectIngredientsForRecipes(recipeIDs []int64) SelectStatement {
	riCols := RecipeIngredients.AllColumns()
	iCols := Ingredients.AllColumns()
	uCols := Units.AllColumns()

	// Build column list
	cols := make([]Projection, 0, len(riCols)+len(iCols)+len(uCols)+1)
	cols = append(cols, riCols...)
	cols = append(cols, iCols...)
	cols = append(cols, uCols...)
	cols = append(cols, RecipeSteps.ID.AS("step_id"))

	return SELECT(cols[0], cols[1:]...).FROM(
		RecipeIngredients.
			INNER_JOIN(RecipeSteps, RecipeIngredients.StepID.EQ(RecipeSteps.ID)).
			INNER_JOIN(Ingredients, RecipeIngredients.IngredientID.EQ(Ingredients.ID)).
			INNER_JOIN(Units, RecipeIngredients.UnitID.EQ(Units.ID)),
	).WHERE(RecipeSteps.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(RecipeIngredients.SortOrder.ASC())
}

// SelectNutrientsForRecipes returns a query to fetch nutrients for recipe ingredients
func SelectNutrientsForRecipes(recipeIDs []int64) SelectStatement {
	inCols := IngredientNutrients.AllColumns()
	nCols := Nutrients.AllColumns()

	// Build column list
	cols := make([]Projection, 0, len(inCols)+len(nCols))
	cols = append(cols, inCols...)
	cols = append(cols, nCols...)

	return SELECT(cols[0], cols[1:]...).FROM(
		IngredientNutrients.
			INNER_JOIN(Nutrients, IngredientNutrients.NutrientID.EQ(Nutrients.ID)).
			INNER_JOIN(RecipeIngredients, IngredientNutrients.IngredientID.EQ(RecipeIngredients.IngredientID)).
			INNER_JOIN(RecipeSteps, RecipeIngredients.StepID.EQ(RecipeSteps.ID)),
	).WHERE(RecipeSteps.RecipeID.IN(IntSliceToExpressions(recipeIDs)...)).
		ORDER_BY(IngredientNutrients.IngredientID.ASC(), Nutrients.Name.ASC())
}

// SelectTags returns a paginated query for tags
func SelectTags(lastID int64, lastName string, limit int64) SelectStatement {
	return SELECT(
		Tags.ID,
		Tags.Name,
	).FROM(Tags).
		WHERE(ROW(Tags.Name, Tags.ID).GT(ROW(String(lastName), Int(lastID)))).
		ORDER_BY(Tags.Name.ASC(), Tags.ID.ASC()).
		LIMIT(limit)
}

// SelectMealPlan returns a query for meal plan entries with recipes
func SelectMealPlan(userID int64, from, until string, lastDate string, lastID, limit int64) SelectStatement {
	return SELECT(
		MealPlan.ID,
		MealPlan.Date,
		MealPlan.UserID,
		MealPlan.RecipeID,
		MealPlan.SortOrder,
		Recipes.ID,
		Recipes.Name,
		Recipes.Servings,
		Recipes.Minutes,
		Recipes.Description,
		Recipes.CreatedBy,
		Recipes.CreatedAt,
	).FROM(
		MealPlan.INNER_JOIN(Recipes, Recipes.ID.EQ(MealPlan.RecipeID)),
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
	return RecipeSteps.INSERT(RecipeSteps.RecipeID, RecipeSteps.Instructions, RecipeSteps.SortOrder).
		VALUES(recipeID, instructions, sortOrder).
		RETURNING(RecipeSteps.AllColumns()...)
}

// InsertStepIngredient returns an insert statement for a step ingredient
func InsertStepIngredient(stepID, ingredientID, unitID int64, amount float64, sortOrder int64) InsertStatement {
	return RecipeIngredients.INSERT(
		RecipeIngredients.StepID,
		RecipeIngredients.IngredientID,
		RecipeIngredients.UnitID,
		RecipeIngredients.Amount,
		RecipeIngredients.SortOrder,
	).VALUES(stepID, ingredientID, unitID, amount, sortOrder)
}

// InsertRecipeImage returns an insert statement for a recipe image
func InsertRecipeImage(recipeID int64, path string, sortOrder int64) InsertStatement {
	return RecipeImages.INSERT(RecipeImages.RecipeID, RecipeImages.Path, RecipeImages.SortOrder).
		VALUES(recipeID, path, sortOrder)
}

// InsertRecipeTag returns an insert statement for a recipe tag
func InsertRecipeTag(recipeID, tagID int64) InsertStatement {
	return RecipeTags.INSERT(RecipeTags.RecipeID, RecipeTags.TagID).
		VALUES(recipeID, tagID)
}

// InsertMealPlan returns an insert statement for a meal plan entry
func InsertMealPlan(date string, userID, recipeID, sortOrder int64) InsertStatement {
	return MealPlan.INSERT(MealPlan.Date, MealPlan.UserID, MealPlan.RecipeID, MealPlan.SortOrder).
		VALUES(date, userID, recipeID, sortOrder)
}

// UpdateRecipe returns an update statement for a recipe
func UpdateRecipe(id int64, name, description string, servings, minutes int64) UpdateStatement {
	return Recipes.UPDATE(Recipes.Name, Recipes.Servings, Recipes.Minutes, Recipes.Description).
		SET(name, servings, minutes, description).
		WHERE(Recipes.ID.EQ(Int(id)))
}

// DeleteRecipeSteps returns a delete statement for recipe steps
func DeleteRecipeSteps(recipeID int64) DeleteStatement {
	return RecipeSteps.DELETE().WHERE(RecipeSteps.RecipeID.EQ(Int(recipeID)))
}

// DeleteRecipeImages returns a delete statement for recipe images
func DeleteRecipeImages(recipeID int64) DeleteStatement {
	return RecipeImages.DELETE().WHERE(RecipeImages.RecipeID.EQ(Int(recipeID)))
}

// DeleteRecipeTags returns a delete statement for recipe tags
func DeleteRecipeTags(recipeID int64) DeleteStatement {
	return RecipeTags.DELETE().WHERE(RecipeTags.RecipeID.EQ(Int(recipeID)))
}

// DeleteRecipeIngredients returns a delete statement for recipe ingredients
func DeleteRecipeIngredients(recipeID int64) DeleteStatement {
	return RecipeIngredients.DELETE().
		WHERE(RecipeIngredients.StepID.IN(
			SELECT(RecipeSteps.ID).
				FROM(RecipeSteps).
				WHERE(RecipeSteps.RecipeID.EQ(Int(recipeID))),
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
