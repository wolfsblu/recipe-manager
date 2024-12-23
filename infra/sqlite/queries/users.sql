-- name: CreatePasswordResetToken :one
INSERT INTO password_resets (user_id, token)
VALUES (?, ?)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (?, ?)
RETURNING *;

-- name: CreateUserRegistration :one
INSERT INTO user_registrations (user_id, token)
VALUES (?, ?)
RETURNING *;

-- name: DeletePasswordResetsBefore :exec
DELETE
FROM password_resets
WHERE created_at < ?;

-- name: DeletePasswordResetTokenByUserId :exec
DELETE
FROM password_resets
WHERE user_id = ?;

-- name: DeleteRegistrationByUserId :exec
DELETE
FROM user_registrations
WHERE user_id = ?;

-- name: DeleteRegistrationsBefore :exec
DELETE
FROM users
WHERE created_at < ?;

-- name: GetPasswordResetToken :one
SELECT sqlc.embed(password_resets), sqlc.embed(users)
FROM password_resets
         INNER JOIN users ON users.id = password_resets.user_id
WHERE token = ?
LIMIT 1;

-- name: GetPasswordResetTokenByUser :one
SELECT *
FROM password_resets
WHERE user_id = ?
LIMIT 1;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;

-- name: GetUserRegistration :one
SELECT sqlc.embed(user_registrations), sqlc.embed(users)
FROM user_registrations
         INNER JOIN users ON users.id = user_registrations.user_id
WHERE token = ?
LIMIT 1;

-- name: UpdatePasswordByUserId :exec
UPDATE users
SET password_hash = ?
WHERE id = ?;

-- name: UpdateUser :exec
UPDATE users
SET email        = ?,
    is_confirmed = ?
WHERE id = ?;