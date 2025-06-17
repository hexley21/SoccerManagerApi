-- name: InsertTeam :one
INSERT INTO teams (id, user_id, name, country_code, budget, total_players) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListTeamsCursor :many
SELECT * FROM teams WHERE id > $1 ORDER BY id LIMIT $2;

-- name: GetTeamByID :one
SELECT * FROM teams WHERE id = $1;

-- name: GetTeamByUserID :one
SELECT * FROM teams WHERE user_id = $1;

-- name: UpdateTeamDataByUserID :exec
UPDATE teams SET name = $2, country_code = $3 WHERE user_id = $1;

-- name: DeleteTeamByID :exec
DELETE FROM teams WHERE id = $1;
