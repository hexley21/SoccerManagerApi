-- name: ListTransfers :many
SELECT * FROM transfers WHERE id > $1 ORDER BY id LIMIT $2;

-- name: ListTransfersByTeamId :many
SELECT * FROM transfers WHERE seller_team_id = $1 AND id > $2 ORDER BY id LIMIT $3;

-- name: GetTransferByPlayerID :one
SELECT * FROM transfers WHERE player_id = $1;

-- name: SelectTransferById :one
SELECT * FROM transfers WHERE id = $1;

-- name: InsertTransferRecordByUser :one
WITH team AS (SELECT id FROM teams WHERE user_id = $1 LIMIT 1)
INSERT INTO transfers (id, player_id, seller_team_id, price)
VALUES ($2, $3, (SELECT id FROM team), $4)
RETURNING id;

-- name: DeleteTransferByID :exec
DELETE FROM transfers WHERE id = $1;

-- name: DeleteTransferByIDAndUserID :exec
DELETE FROM transfers USING teams WHERE transfers.id = $1 AND transfers.seller_team_id = teams.id AND teams.user_id = $2;

-- name: UpdateTransferPriceByIDAndUserID :exec
UPDATE transfers SET price = $2 FROM teams WHERE transfers.id = $1 AND transfers.seller_team_id = teams.id AND teams.user_id = $3;
