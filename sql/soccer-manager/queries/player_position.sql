-- name: GetAllPositionCodes :many
SELECT code FROM positions ORDER BY code;

-- name: GetAllPositions :many
SELECT * FROM positions ORDER BY code;

-- name: GetPositionTranslationsByLocale :many
SELECT * FROM position_translations WHERE locale = $1 ORDER BY locale;

-- name: GetPositionTranslationsByPositionCode :many
SELECT * FROM position_translations WHERE position_code = $1 ORDER BY locale;

-- name: GetPositionTranslations :many
SELECT * FROM position_translations ORDER BY position_code;

-- name: InsertPositionTranslation :exec
INSERT INTO position_translations (position_code, locale, label) VALUES ($1, $2, $3);

-- name: UpdatePositionTranslationLabel :exec
UPDATE position_translations SET label = $3 WHERE position_code = $1 AND locale = $2;

-- name: DeletePositionTranslation :exec
DELETE FROM position_translations WHERE position_code = $1 AND locale = $2;
