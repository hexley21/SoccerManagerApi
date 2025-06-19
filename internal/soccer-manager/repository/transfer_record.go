package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/hexley21/soccer-manager/pkg/infra/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mock/mock_transfer_record.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository TransferRecordRepository
type TransferRecordRepository interface {
	ListTransferRecords(ctx context.Context, arg ListTransferRecordsParams) ([]TransferRecord, error)
	GetTransferRecordByID(ctx context.Context, id int64) (TransferRecord, error)
}

type pgTransferRecordRepository struct {
	db            *pgxpool.Pool
	snowflakeNode *snowflake.Node
}

func NewTransferRecordRepository(
	db *pgxpool.Pool,
	snowflakeNode *snowflake.Node,
) *pgTransferRecordRepository {
	return &pgTransferRecordRepository{db: db, snowflakeNode: snowflakeNode}
}

const insertTransferRecord = `-- name: InsertTransferRecord :exec
INSERT INTO transfer_records (id, player_id, seller_team_id, buyer_team_id, sold_price, listed_at) VALUES ($1, $2, $3, $4, $5, $6)
`

type InsertTransferRecordParams struct {
	ID           int64              `json:"id"`
	PlayerID     int64              `json:"player_id"`
	SellerTeamID int64              `json:"seller_team_id"`
	BuyerTeamID  int64              `json:"buyer_team_id"`
	SoldPrice    int64              `json:"sold_price"`
	ListedAt     pgtype.Timestamptz `json:"listed_at"`
}

func insertTransferRecordWithQuerier(
	ctx context.Context,
	querier postgres.Querier,
	arg InsertTransferRecordParams,
) error {
	_, err := querier.Exec(ctx, insertTransferRecord,
		arg.ID,
		arg.PlayerID,
		arg.SellerTeamID,
		arg.BuyerTeamID,
		arg.SoldPrice,
		arg.ListedAt,
	)
	return err
}

const getTransferRecordByID = `-- name: GetTransferRecordByID :one
SELECT id, player_id, seller_team_id, buyer_team_id, sold_price, listed_at, sold_at FROM transfer_records WHERE id = $1
`

func (r *pgTransferRecordRepository) GetTransferRecordByID(
	ctx context.Context,
	id int64,
) (TransferRecord, error) {
	row := r.db.QueryRow(ctx, getTransferRecordByID, id)
	var i TransferRecord
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.SellerTeamID,
		&i.BuyerTeamID,
		&i.SoldPrice,
		&i.ListedAt,
		&i.SoldAt,
	)
	return i, err
}

const listTransferRecords = `-- name: ListTransferRecords :many
SELECT id, player_id, seller_team_id, buyer_team_id, sold_price, listed_at, sold_at FROM transfer_records WHERE id > $1 ORDER BY id LIMIT $2
`

type ListTransferRecordsParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

func (r *pgTransferRecordRepository) ListTransferRecords(
	ctx context.Context,
	arg ListTransferRecordsParams,
) ([]TransferRecord, error) {
	rows, err := r.db.Query(ctx, listTransferRecords, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TransferRecord{}
	for rows.Next() {
		var i TransferRecord
		if err := rows.Scan(
			&i.ID,
			&i.PlayerID,
			&i.SellerTeamID,
			&i.BuyerTeamID,
			&i.SoldPrice,
			&i.ListedAt,
			&i.SoldAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
