-- name: BrowseRecipes :many
SELECT *
FROM recipes;

-- name: CreateRecipe :one
INSERT INTO recipes (name, servings, minutes, description, created_by)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: CreateRecipeImages :one
INSERT INTO recipe_images (recipe_id, path, sort_order)
VALUES (?, ?, ?)
RETURNING id;

-- name: CreateRecipeStep :one
INSERT INTO recipe_steps (recipe_id, instructions, sort_order)
VALUES (?, ?, ?)
RETURNING id;

-- name: CreateStepIngredient :one
INSERT INTO recipe_ingredients (step_id, ingredient_id, unit_id, amount, sort_order)
VALUES (?, ?, ?, ?, ?)
RETURNING id;

-- name: CreateRecipeTag :exec
INSERT INTO recipe_tags (recipe_id, tag_id)
VALUES (?, ?);

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
       units.symbol as unit_symbol, 
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

-- name: GetTags :many
SELECT *
FROM tags
ORDER BY name;

-- name: GetStepsForRecipes :many
SELECT id, instructions, sort_order, recipe_id
FROM recipe_steps
WHERE recipe_id IN (
    sqlc.slice(recipe_ids)
    )
ORDER BY sort_order;

-- name: GetTagsForRecipes :many
SELECT recipe_tags.recipe_id, sqlc.embed(tags)
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
SET name = ?, servings = ?, minutes = ?, description = ?
WHERE id = ?;

-- name: DeleteRecipeIngredients :exec
DELETE FROM recipe_ingredients
WHERE step_id IN (
    SELECT id FROM recipe_steps WHERE recipe_id = ?
);

-- name: DeleteRecipeSteps :exec
DELETE FROM recipe_steps
WHERE recipe_id = ?;

-- name: DeleteRecipeImages :exec
DELETE FROM recipe_images
WHERE recipe_id = ?;

-- name: DeleteRecipeTags :exec
DELETE FROM recipe_tags
WHERE recipe_id = ?;

-- name: AddVote :exec
INSERT INTO recipe_votes (recipe_id, user_id, vote)
VALUES (?, ?, ?)
ON CONFLICT (recipe_id, user_id) DO UPDATE SET vote = excluded.vote;

-- name: RemoveVote :exec
DELETE FROM recipe_votes
WHERE recipe_id = ? AND user_id = ?;

-- name: GetRecipeVotes :one
SELECT
    COALESCE(SUM(vote), 0) as total_score
FROM recipe_votes
WHERE recipe_id = ?;

-- name: GetUserVote :one
SELECT vote
FROM recipe_votes
WHERE recipe_id = ? AND user_id = ?
LIMIT 1;

-- name: GetVotesForRecipes :many
SELECT
    recipe_id,
    COALESCE(SUM(vote), 0) as total_score
FROM recipe_votes
WHERE recipe_id IN (sqlc.slice(recipe_ids))
GROUP BY recipe_id;

-- name: GetUserVotesForRecipes :many
SELECT
    recipe_id,
    vote
FROM recipe_votes
WHERE recipe_id IN (sqlc.slice(recipe_ids)) AND user_id = ?;

-- name: CreateIngredient :one
INSERT INTO ingredients (name)
VALUES (?)
RETURNING id;

-- name: UpdateIngredient :exec
UPDATE ingredients
SET name = ?
WHERE id = ?;

-- name: DeleteIngredient :exec
DELETE FROM ingredients
WHERE id = ?;

-- name: CreateUnit :one
INSERT INTO units (name, symbol)
VALUES (?, ?)
RETURNING id;

-- name: UpdateUnit :exec
UPDATE units
SET name = ?, symbol = ?
WHERE id = ?;

-- name: DeleteUnit :exec
DELETE FROM units
WHERE id = ?;

-- name: CreateMealPlan :exec
INSERT INTO meal_plan (date, user_id, recipe_id, sort_order)
VALUES (?, ?, ?, ?);

-- name: DeleteMealPlan :exec
DELETE FROM meal_plan
WHERE user_id = ? AND recipe_id = ? AND date = ?;

-- name: GetNutrients :many
SELECT *
FROM nutrients
ORDER BY name;

-- name: CreateNutrient :one
INSERT INTO nutrients (name, unit)
VALUES (?, ?)
RETURNING id;

-- name: UpdateNutrient :exec
UPDATE nutrients
SET name = ?, unit = ?
WHERE id = ?;

-- name: DeleteNutrient :exec
DELETE FROM nutrients
WHERE id = ?;

-- name: GetNutrientsForIngredient :many
SELECT sqlc.embed(nutrients), ingredient_nutrients.amount
FROM ingredient_nutrients
INNER JOIN nutrients ON ingredient_nutrients.nutrient_id = nutrients.id
WHERE ingredient_nutrients.ingredient_id = ?
ORDER BY nutrients.name;

-- name: GetNutrientsForIngredients :many
SELECT ingredient_nutrients.ingredient_id, sqlc.embed(nutrients), ingredient_nutrients.amount
FROM ingredient_nutrients
INNER JOIN nutrients ON ingredient_nutrients.nutrient_id = nutrients.id
WHERE ingredient_nutrients.ingredient_id IN (sqlc.slice(ingredient_ids))
ORDER BY ingredient_nutrients.ingredient_id, nutrients.name;

-- name: GetNutrientsForRecipes :many
SELECT ingredient_nutrients.ingredient_id, sqlc.embed(nutrients), ingredient_nutrients.amount
FROM ingredient_nutrients
INNER JOIN nutrients ON ingredient_nutrients.nutrient_id = nutrients.id
INNER JOIN recipe_ingredients ON ingredient_nutrients.ingredient_id = recipe_ingredients.ingredient_id
INNER JOIN recipe_steps ON recipe_ingredients.step_id = recipe_steps.id
WHERE recipe_steps.recipe_id IN (sqlc.slice(recipe_ids))
ORDER BY ingredient_nutrients.ingredient_id, nutrients.name;

-- name: AddIngredientNutrient :exec
INSERT INTO ingredient_nutrients (ingredient_id, nutrient_id, amount)
VALUES (?, ?, ?)
ON CONFLICT (ingredient_id, nutrient_id) DO UPDATE SET amount = excluded.amount;

-- name: DeleteIngredientNutrients :exec
DELETE FROM ingredient_nutrients
WHERE ingredient_id = ?;