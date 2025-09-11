-- name: BrowseRecipes :many
SELECT *
FROM recipes;

-- name: CreateRecipe :one
INSERT INTO recipes (name, servings, minutes, description, created_by)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: CreateRecipeStep :one
INSERT INTO recipe_steps (recipe_id, instructions)
VALUES (?, ?)
RETURNING id;

-- name: CreateStepIngredient :one
INSERT INTO recipe_ingredients (step_id, ingredient_id, unit_id, amount, sort_order)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

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

-- name: GetImagesForRecipes :many
SELECT id, path, sort_order, recipe_id
FROM recipe_images
WHERE recipe_id IN (
    sqlc.slice(recipe_ids)
    )
ORDER BY sort_order;

-- name: GetIngredientsForRecipes :many
SELECT recipe_ingredients.id as recipe_ingredient_id, 
       ingredients.id as ingredient_id, 
       ingredients.name as ingredient_name, 
       units.id as unit_id, 
       units.name as unit_name, 
       units.code as unit_code, 
       recipe_steps.id as step_id, 
       recipe_ingredients.amount, 
       recipe_ingredients.sort_order
FROM recipe_ingredients
         INNER JOIN recipe_steps ON recipe_ingredients.step_id = recipe_steps.id
         INNER JOIN ingredients ON recipe_ingredients.ingredient_id = ingredients.id
         INNER JOIN units ON recipe_ingredients.unit_id = units.id
WHERE recipe_steps.recipe_id IN (
    sqlc.slice(recipe_ids)
    )
ORDER BY recipe_ingredients.sort_order;

-- name: GetIngredients :many
SELECT *
FROM ingredients
ORDER BY name;

-- name: GetUnits :many
SELECT *
FROM units
ORDER BY name;

-- name: GetStepsForRecipes :many
SELECT id, instructions, sort_order, recipe_id
FROM recipe_steps
WHERE recipe_id IN (
    sqlc.slice(recipe_ids)
    )
ORDER BY sort_order;

-- name: GetTagsForRecipes :many
SELECT tags.id, tags.name, recipe_tags.recipe_id
FROM tags
         INNER JOIN recipe_tags ON tags.id = recipe_tags.tag_id
WHERE recipe_tags.recipe_id IN (
    sqlc.slice(recipe_ids)
    )
ORDER BY tags.name;

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