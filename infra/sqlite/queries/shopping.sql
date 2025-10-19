-- name: GetShoppingListsByUserID :many
SELECT id, user_id, name FROM shopping_lists
WHERE user_id = sqlc.arg(user_id)
  AND id > COALESCE(sqlc.narg(cursor), 0)
ORDER BY id ASC
LIMIT sqlc.arg(limit);

-- name: GetShoppingListByID :one
SELECT id, user_id, name FROM shopping_lists
WHERE id = ?;

-- name: CreateShoppingList :one
INSERT INTO shopping_lists (user_id, name)
VALUES (?, ?)
RETURNING id, user_id, name;

-- name: UpdateShoppingList :one
UPDATE shopping_lists
SET name = ?
WHERE id = ?
RETURNING id, user_id, name;

-- name: DeleteShoppingList :exec
DELETE FROM shopping_lists
WHERE id = ?;

-- name: GetShoppingListItemsByListID :many
SELECT id, shopping_list_id, ingredient, quantity, unit, done, sort_order
FROM shopping_list_items
WHERE shopping_list_id = ?
ORDER BY sort_order ASC;

-- name: GetShoppingListItemByID :one
SELECT id, shopping_list_id, ingredient, quantity, unit, done, sort_order
FROM shopping_list_items
WHERE id = ? AND shopping_list_id = ?;

-- name: CreateShoppingListItem :one
INSERT INTO shopping_list_items (shopping_list_id, ingredient, quantity, unit, done, sort_order)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, shopping_list_id, ingredient, quantity, unit, done, sort_order;

-- name: UpdateShoppingListItem :one
UPDATE shopping_list_items
SET ingredient = ?, quantity = ?, unit = ?, done = ?
WHERE id = ?
RETURNING id, shopping_list_id, ingredient, quantity, unit, done, sort_order;

-- name: DeleteShoppingListItem :exec
DELETE FROM shopping_list_items
WHERE id = ?;

-- name: UpdateShoppingListItemSortOrder :exec
UPDATE shopping_list_items
SET sort_order = ?
WHERE id = ?;