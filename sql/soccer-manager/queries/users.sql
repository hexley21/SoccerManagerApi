-- name: GetUserByID :one
SELECT id, username, role FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, role FROM users WHERE username = $1;

-- name: ListUsersCursor :many
SELECT id, username, role FROM users WHERE id > $1 ORDER BY id LIMIT $2;

-- name: GetUserHashByID :one
SELECT hash FROM users WHERE id = $1 LIMIT 1;

-- name: GetAuth :one
SELECT id, hash, role FROM users WHERE username = $1 LIMIT 1;

-- name: CheckUserExists :one
SELECT EXISTS (SELECT 1 FROM users WHERE id = $1) AS user_exists;

-- name: CreateUser :one
INSERT INTO users (id, username, role, hash) VALUES ($1, $2, $3, $4) RETURNING id, username, role;

-- name: UpdateUserHash :exec
UPDATE users SET hash = $2 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
