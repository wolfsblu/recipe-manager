package queries

import (
	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectIngredientsPaginated returns a query for paginated ingredients with sorting and filtering
func SelectIngredientsPaginated(lastID int32, lastName string, limit int64, sortField, sortOrder, search string) SelectStatement {
	var comparison BoolExpression
	var orderBy []OrderByClause
	var whereClauses []BoolExpression

	// Build dynamic WHERE clause based on sort field and order
	switch sortField {
	case "id":
		if sortOrder == "asc" {
			comparison = Ingredient.ID.GT(Int32(lastID))
			orderBy = []OrderByClause{Ingredient.ID.ASC()}
		} else {
			comparison = Ingredient.ID.LT(Int32(lastID))
			orderBy = []OrderByClause{Ingredient.ID.DESC()}
		}
	default: // name
		if sortOrder == "asc" {
			comparison = ROW(Ingredient.Name, Ingredient.ID).GT(ROW(String(lastName), Int32(lastID)))
			orderBy = []OrderByClause{Ingredient.Name.ASC(), Ingredient.ID.ASC()}
		} else {
			comparison = ROW(Ingredient.Name, Ingredient.ID).LT(ROW(String(lastName), Int32(lastID)))
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
func SelectNutrientsForIngredients(ingredientIDs []int32) SelectStatement {
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
func InsertIngredientNutrient(ingredientID, nutrientID int32, amount float64) InsertStatement {
	return IngredientNutrient.INSERT(
		IngredientNutrient.IngredientID,
		IngredientNutrient.NutrientID,
		IngredientNutrient.Amount,
	).VALUES(ingredientID, nutrientID, amount)
}

// UpdateIngredient returns an update statement for an ingredient
func UpdateIngredient(id int32, name string) UpdateStatement {
	return Ingredient.UPDATE(Ingredient.Name).
		SET(name).
		WHERE(Ingredient.ID.EQ(Int32(id)))
}

// DeleteIngredient returns a delete statement for an ingredient
func DeleteIngredient(id int32) DeleteStatement {
	return Ingredient.DELETE().WHERE(Ingredient.ID.EQ(Int32(id)))
}

// DeleteIngredientNutrients returns a delete statement for ingredient nutrients
func DeleteIngredientNutrients(ingredientID int32) DeleteStatement {
	return IngredientNutrient.DELETE().
		WHERE(IngredientNutrient.IngredientID.EQ(Int32(ingredientID)))
}

// SelectUnits returns a paginated query for units
func SelectUnits(lastID int32, lastName string, limit int32) SelectStatement {
	return SELECT(
		Unit.ID,
		Unit.Name,
		Unit.Symbol,
	).FROM(Unit).
		WHERE(ROW(Unit.Name, Unit.ID).GT(ROW(String(lastName), Int32(lastID)))).
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
func UpdateUnit(id int32, name string, symbol *string) UpdateStatement {
	return Unit.UPDATE(Unit.Name, Unit.Symbol).
		SET(name, symbol).
		WHERE(Unit.ID.EQ(Int32(id)))
}

// DeleteUnit returns a delete statement for a unit
func DeleteUnit(id int32) DeleteStatement {
	return Unit.DELETE().WHERE(Unit.ID.EQ(Int32(id)))
}
