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
			comparison = Ingredients.ID.GT(Int(lastID))
			orderBy = []OrderByClause{Ingredients.ID.ASC()}
		} else {
			comparison = Ingredients.ID.LT(Int(lastID))
			orderBy = []OrderByClause{Ingredients.ID.DESC()}
		}
	default: // name
		if sortOrder == "asc" {
			comparison = ROW(Ingredients.Name, Ingredients.ID).GT(ROW(String(lastName), Int(lastID)))
			orderBy = []OrderByClause{Ingredients.Name.ASC(), Ingredients.ID.ASC()}
		} else {
			comparison = ROW(Ingredients.Name, Ingredients.ID).LT(ROW(String(lastName), Int(lastID)))
			orderBy = []OrderByClause{Ingredients.Name.DESC(), Ingredients.ID.DESC()}
		}
	}

	whereClauses = append(whereClauses, comparison)

	// Add search filter if provided
	if search != "" {
		whereClauses = append(whereClauses, Ingredients.Name.LIKE(String("%"+search+"%")))
	}

	// Combine WHERE clauses with AND
	finalWhere := whereClauses[0]
	for i := 1; i < len(whereClauses); i++ {
		finalWhere = finalWhere.AND(whereClauses[i])
	}

	return SELECT(
		Ingredients.ID,
		Ingredients.Name,
	).FROM(Ingredients).
		WHERE(finalWhere).
		ORDER_BY(orderBy...).
		LIMIT(limit)
}

// SelectNutrientsForIngredients returns a query to fetch nutrients for given ingredient IDs
func SelectNutrientsForIngredients(ingredientIDs []int64) SelectStatement {
	return SELECT(
		IngredientNutrients.IngredientID,
		Nutrients.ID,
		Nutrients.Name,
		Nutrients.Unit,
		IngredientNutrients.Amount,
	).FROM(
		IngredientNutrients.INNER_JOIN(Nutrients, IngredientNutrients.NutrientID.EQ(Nutrients.ID)),
	).WHERE(IngredientNutrients.IngredientID.IN(IntSliceToExpressions(ingredientIDs)...))
}

// InsertIngredient returns an insert statement for a new ingredient
func InsertIngredient(name string) InsertStatement {
	return Ingredients.INSERT(Ingredients.Name).
		VALUES(name).
		RETURNING(Ingredients.AllColumns()...)
}

// InsertIngredientNutrient returns an insert statement for an ingredient nutrient
func InsertIngredientNutrient(ingredientID, nutrientID int64, amount float64) InsertStatement {
	return IngredientNutrients.INSERT(
		IngredientNutrients.IngredientID,
		IngredientNutrients.NutrientID,
		IngredientNutrients.Amount,
	).VALUES(ingredientID, nutrientID, amount)
}

// UpdateIngredient returns an update statement for an ingredient
func UpdateIngredient(id int64, name string) UpdateStatement {
	return Ingredients.UPDATE(Ingredients.Name).
		SET(name).
		WHERE(Ingredients.ID.EQ(Int(id)))
}

// DeleteIngredient returns a delete statement for an ingredient
func DeleteIngredient(id int64) DeleteStatement {
	return Ingredients.DELETE().WHERE(Ingredients.ID.EQ(Int(id)))
}

// DeleteIngredientNutrients returns a delete statement for ingredient nutrients
func DeleteIngredientNutrients(ingredientID int64) DeleteStatement {
	return IngredientNutrients.DELETE().
		WHERE(IngredientNutrients.IngredientID.EQ(Int(ingredientID)))
}

// SelectUnits returns a paginated query for units
func SelectUnits(lastID int64, lastName string, limit int64) SelectStatement {
	return SELECT(
		Units.ID,
		Units.Name,
		Units.Symbol,
	).FROM(Units).
		WHERE(ROW(Units.Name, Units.ID).GT(ROW(String(lastName), Int(lastID)))).
		ORDER_BY(Units.Name.ASC(), Units.ID.ASC()).
		LIMIT(limit)
}

// InsertUnit returns an insert statement for a new unit
func InsertUnit(name string, symbol *string) InsertStatement {
	return Units.INSERT(Units.Name, Units.Symbol).
		VALUES(name, symbol).
		RETURNING(Units.AllColumns()...)
}

// UpdateUnit returns an update statement for a unit
func UpdateUnit(id int64, name string, symbol *string) UpdateStatement {
	return Units.UPDATE(Units.Name, Units.Symbol).
		SET(name, symbol).
		WHERE(Units.ID.EQ(Int(id)))
}

// DeleteUnit returns a delete statement for a unit
func DeleteUnit(id int64) DeleteStatement {
	return Units.DELETE().WHERE(Units.ID.EQ(Int(id)))
}
