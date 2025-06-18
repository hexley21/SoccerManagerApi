-- name: ListPlayersByCursor :many
SELECT id, team_id, country_code, first_name, last_name, age, position_code, price FROM players WHERE id > $1 ORDER BY id LIMIT $2;

-- name: GetPlayerByID :one
SELECT id, team_id, country_code, first_name, last_name, age, position_code, price FROM players WHERE id = $1;

-- name: UpdatePlayerNameAndCountry :exec
UPDATE players SET first_name = $3, last_name = $4, country_code = $5 FROM teams WHERE players.team_id = teams.id AND players.id = $2 AND teams.user_id = $1;

-- name: UpdatePlayerPriceAndTeam :exec
UPDATE players SET price = $2, team_id = $3 WHERE id = $1;

-- name: ListPlayersByUserID :many
SELECT p.* FROM players p JOIN teams t ON p.team_id = t.id WHERE t.user_id = $1 AND p.id > $2 ORDER BY p.id LIMIT $3;

-- name: ListPlayersByTeamID :many
SELECT p.* FROM players p WHERE p.team_id = $1 AND p.id > $2 ORDER BY p.id LIMIT $3;

-- name: InsertPlayer :exec
INSERT INTO players (id, team_id, country_code, first_name, last_name, age, position_code, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
