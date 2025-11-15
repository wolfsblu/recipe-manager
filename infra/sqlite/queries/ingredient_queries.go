package queries

import (
	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectIngredientsPaginated returns a query for paginated ingredients with sorting and filtering
func SelectIngredientsPaginated(lastID int64, lastName string, limit int64, sortField, sortOrder, search string) SelectStatement {
	var comparison BoolExpression
	var orderBy []OrderByClause
	var whereClauses []BoolExpression

	// Build dynamic WHERE clause based on sort field and order
	switch sortField {
	case "id":
		if sortOrder == "asc" {
			comparison = Ingredient.ID.GT(Int(lastID))
			orderBy = []OrderByClause{Ingredient.ID.ASC()}
		} else {
			comparison = Ingredient.ID.LT(Int(lastID))
			orderBy = []OrderByClause{Ingredient.ID.DESC()}
		}
	default: // name
		if sortOrder == "asc" {
			comparison = ROW(Ingredient.Name, Ingredient.ID).GT(ROW(String(lastName), Int(lastID)))
			orderBy = []OrderByClause{Ingredient.Name.ASC(), Ingredient.ID.ASC()}
		} else {
			comparison = ROW(Ingredient.Name, Ingredient.ID).LT(ROW(String(lastName), Int(lastID)))
			orderBy = []OrderByClause{Ingredient.Name.DESC(), Ingredient.ID.DESC()}
		}
	}

	whereClauses = append(whereClauses, comparison)

	// Add search filter if provided
	if search != "" {
		whereClauses = append(whereClauses, Ingredient.Name.LIKE(String("%"+search+"%")))
	}

	// Combine WHERE clauses with AND
	finalWhere := whereClauses[0]
	for i := 1; i < len(whereClauses); i++ {
		finalWhere = finalWhere.AND(whereClauses[i])
	}

	return SELECT(
		Ingredient.ID,
		Ingredient.Name,
	).FROM(Ingredient).
		WHERE(finalWhere).
		ORDER_BY(orderBy...).
		LIMIT(limit)
}

// SelectNutrientsForIngredients returns a query to fetch nutrients for given ingredient IDs
func SelectNutrientsForIngredients(ingredientIDs []int64) SelectStatement {
	return SELECT(
		IngredientNutrient.IngredientID,
		Nutrient.ID,
		Nutrient.Name,
		Nutrient.Unit,
		IngredientNutrient.Amount,
	).FROM(
		IngredientNutrient.INNER_JOIN(Nutrient, IngredientNutrient.NutrientID.EQ(Nutrient.ID)),
	).WHERE(IngredientNutrient.IngredientID.IN(IntSliceToExpressions(ingredientIDs)...))
}

// InsertIngredient returns an insert statement for a new ingredient
func InsertIngredient(name string) InsertStatement {
	return Ingredient.INSERT(Ingredient.Name).
		VALUES(name).
		RETURNING(Ingredient.AllColumns)
}

// InsertIngredientNutrient returns an insert statement for an ingredient nutrient
func InsertIngredientNutrient(ingredientID, nutrientID int64, amount float64) InsertStatement {
	return IngredientNutrient.INSERT(
		IngredientNutrient.IngredientID,
		IngredientNutrient.NutrientID,
		IngredientNutrient.Amount,
	).VALUES(ingredientID, nutrientID, amount)
}

// UpdateIngredient returns an update statement for an ingredient
func UpdateIngredient(id int64, name string) UpdateStatement {
	return Ingredient.UPDATE(Ingredient.Name).
		SET(name).
		WHERE(Ingredient.ID.EQ(Int(id)))
}

// DeleteIngredient returns a delete statement for an ingredient
func DeleteIngredient(id int64) DeleteStatement {
	return Ingredient.DELETE().WHERE(Ingredient.ID.EQ(Int(id)))
}

// DeleteIngredientNutrients returns a delete statement for ingredient nutrients
func DeleteIngredientNutrients(ingredientID int64) DeleteStatement {
	return IngredientNutrient.DELETE().
		WHERE(IngredientNutrient.IngredientID.EQ(Int(ingredientID)))
}

// SelectUnits returns a paginated query for units
func SelectUnits(lastID int64, lastName string, limit int64) SelectStatement {
	return SELECT(
		Unit.ID,
		Unit.Name,
		Unit.Symbol,
	).FROM(Unit).
		WHERE(ROW(Unit.Name, Unit.ID).GT(ROW(String(lastName), Int(lastID)))).
		ORDER_BY(Unit.Name.ASC(), Unit.ID.ASC()).
		LIMIT(limit)
}

// InsertUnit returns an insert statement for a new unit
func InsertUnit(name string, symbol *string) InsertStatement {
	return Unit.INSERT(Unit.Name, Unit.Symbol).
		VALUES(name, symbol).
		RETURNING(Unit.AllColumns)
}

// UpdateUnit returns an update statement for a unit
func UpdateUnit(id int64, name string, symbol *string) UpdateStatement {
	return Unit.UPDATE(Unit.Name, Unit.Symbol).
		SET(name, symbol).
		WHERE(Unit.ID.EQ(Int(id)))
}

// DeleteUnit returns a delete statement for a unit
func DeleteUnit(id int64) DeleteStatement {
	return Unit.DELETE().WHERE(Unit.ID.EQ(Int(id)))
}
