-- name: GetTranslationsByUserID :many
SELECT tt.team_id, tt.locale, tt.name FROM team_translations tt JOIN teams t ON t.id = tt.team_id WHERE t.user_id = $1;

-- name: InsertTranslationByUserID :exec
INSERT INTO team_translations (team_id, locale, name) SELECT t.id, $2, $3 FROM teams t WHERE t.user_id = $1;

-- name: UpdateTranslationNameByUserID :exec
UPDATE team_translations AS tt SET name = $3 FROM teams AS t WHERE tt.team_id = t.id AND t.user_id = $1 AND tt.locale = $2;

-- name: ListLocalesByTeamID :many
SELECT locale FROM team_translations WHERE team_id = $1;

-- name: DeleteTranslationByUserID :exec
DELETE FROM team_translations AS tt USING teams AS t WHERE tt.team_id = t.id AND t.user_id = $1 AND tt.locale = $2;

-- name: ListTranslatedTeamsCursor :many
SELECT 
    t.id,
    t.user_id,
    COALESCE(tt.name, t.name) AS name,
    t.country_code,
    t.budget,
    t.total_players
FROM teams t
LEFT JOIN team_translations tt
    ON t.id = tt.team_id
    AND tt.locale = $2
WHERE t.id > $1
ORDER BY t.id
LIMIT $3;

-- name: GetTranslatedTeamWithId :one
SELECT 
    t.id,
    t.user_id,
    COALESCE(tt.name, t.name) AS name,
    t.country_code,
    t.budget,
    t.total_players
FROM teams t
LEFT JOIN team_translations tt
    ON t.id = tt.team_id
    AND tt.locale = $2
WHERE t.id = $1;

-- name: GetTranslatedTeamWithUserId :one
SELECT 
    t.id,
    t.user_id,
    COALESCE(tt.name, t.name) AS name,
    t.country_code,
    t.budget,
    t.total_players
FROM teams t
LEFT JOIN team_translations tt
    ON t.id = tt.team_id
    AND tt.locale = $2
WHERE t.user_id = $1;