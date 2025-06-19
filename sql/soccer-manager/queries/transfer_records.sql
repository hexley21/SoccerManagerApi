-- name: GetTransferRecordByID :one
SELECT * FROM transfer_records WHERE id = $1;

-- name: ListTransferRecords :many
SELECT * FROM transfer_records WHERE id > $1 ORDER BY id LIMIT $2;

-- name: InsertTransferRecord :exec
INSERT INTO transfer_records (id, player_id, seller_team_id, buyer_team_id, sold_price, listed_at) VALUES ($1, $2, $3, $4, $5, $6);
