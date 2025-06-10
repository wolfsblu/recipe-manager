-- name: BrowseRecipes :many
SELECT *
FROM recipes;

-- name: CreateRecipe :one
INSERT INTO recipes (name, servings, minutes, description, created_by)
VALUES (?, ?, ?, ?, ?)
    RETURNING *;

-- name: DeleteRecipe :exec
DELETE
FROM recipes
WHERE id = ?;

-- name: GetRecipe :one
SELECT *
FROM recipes
WHERE id = ?
LIMIT 1;

-- name: GetMealPlan :many
SELECT sqlc.embed(meal_plan),
       sqlc.embed(recipes)
FROM meal_plan
    INNER JOIN recipes ON meal_plan.recipe_id = recipes.id
WHERE user_id = ?
    AND meal_plan.date >= sqlc.arg(from_date)
    AND meal_plan.date <= sqlc.arg(until_date)
ORDER BY meal_plan.date, meal_plan.sort_order;

-- name: ListRecipes :many
SELECT *
FROM recipes
WHERE created_by = ?
ORDER BY name;

-- name: UpdateRecipe :exec
UPDATE recipes
set name = ?
WHERE id = ?
RETURNING *;