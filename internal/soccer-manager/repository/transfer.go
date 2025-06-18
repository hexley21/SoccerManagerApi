package repository

import (
	"context"
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/hexley21/soccer-manager/pkg/infra/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mock/mock_transfer.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository TransferRepository
type TransferRepository interface {
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	ListTransfersByTeamId(ctx context.Context, arg ListTransfersByTeamIdParams) ([]Transfer, error)
	GetTransferByPlayerID(ctx context.Context, playerID int64) (Transfer, error)
	SelectTransferById(ctx context.Context, id int64) (Transfer, error)

	InsertTransferRecordByUser(
		ctx context.Context,
		arg InsertTransferRecordByUserParams,
	) (int64, error)
	DeleteTransferByIDAndUserID(ctx context.Context, arg DeleteTransferByIDAndUserIDParams) error
	UpdateTransferPriceByIDAndUserID(ctx context.Context, arg UpdateTransferPriceByIDAndUserIDParams) error

	BuyPlayer(ctx context.Context, transferId int64, buyerTeamId int64) error
}

type pgTransferRepository struct {
	db            *pgxpool.Pool
	snowflakeNode *snowflake.Node
}

func NewTransferRepository(db *pgxpool.Pool, snowflakeNode *snowflake.Node) *pgTransferRepository {
	return &pgTransferRepository{
		db:            db,
		snowflakeNode: snowflakeNode,
	}
}

const insertTransferRecordByUser = `-- name: InsertTransferRecordByUser :one
WITH team AS (SELECT id FROM teams WHERE user_id = $1 LIMIT 1)
INSERT INTO transfers (id, player_id, seller_team_id, price)
VALUES ($2, $3, (SELECT id FROM team), $4)
RETURNING id
`

type InsertTransferRecordByUserParams struct {
	UserID   int64 `json:"user_id"`
	PlayerID int64 `json:"player_id"`
	Price    int64 `json:"price"`
}

func (r *pgTransferRepository) InsertTransferRecordByUser(
	ctx context.Context,
	arg InsertTransferRecordByUserParams,
) (int64, error) {
	row := r.db.QueryRow(ctx, insertTransferRecordByUser,
		arg.UserID,
		r.snowflakeNode.Generate().Int64(),
		arg.PlayerID,
		arg.Price,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteTransferByIDAndUserID = `-- name: DeleteTransferByIDAndUserID :exec
DELETE FROM transfers USING teams WHERE transfers.id = $1 AND transfers.seller_team_id = teams.id AND teams.user_id = $2
`

type DeleteTransferByIDAndUserIDParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (r *pgTransferRepository) DeleteTransferByIDAndUserID(ctx context.Context, arg DeleteTransferByIDAndUserIDParams) error {
	res, err := r.db.Exec(ctx, deleteTransferByIDAndUserID, arg.ID, arg.UserID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const deleteTransferByID = `-- name: DeleteTransferByID :exec
DELETE FROM transfers WHERE id = $1
`

func deleteTransferByIDWithQuerier(ctx context.Context, querier postgres.Querier, id int64) error {
	res, err := querier.Exec(ctx, deleteTransferByID, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const selectTransferById = `-- name: SelectTransferById :one
SELECT id, player_id, seller_team_id, price, listed_at FROM transfers WHERE id = $1
`

func (r *pgTransferRepository) SelectTransferById(ctx context.Context, id int64) (Transfer, error) {
	return selectTransferByIdWithQuerier(ctx, r.db, id)
}

func selectTransferByIdWithQuerier(
	ctx context.Context,
	querier postgres.Querier,
	id int64,
) (Transfer, error) {
	row := querier.QueryRow(ctx, selectTransferById, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.SellerTeamID,
		&i.Price,
		&i.ListedAt,
	)
	return i, err
}

const getTransferByPlayerID = `-- name: GetTransferByPlayerID :one
SELECT id, player_id, seller_team_id, price, listed_at FROM transfers WHERE player_id = $1
`

func (r *pgTransferRepository) GetTransferByPlayerID(
	ctx context.Context,
	playerID int64,
) (Transfer, error) {
	row := r.db.QueryRow(ctx, getTransferByPlayerID, playerID)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.SellerTeamID,
		&i.Price,
		&i.ListedAt,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, player_id, seller_team_id, price, listed_at FROM transfers WHERE id > $1 ORDER BY id LIMIT $2
`

type ListTransfersParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

func (r *pgTransferRepository) ListTransfers(
	ctx context.Context,
	arg ListTransfersParams,
) ([]Transfer, error) {
	rows, err := r.db.Query(ctx, listTransfers, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.PlayerID,
			&i.SellerTeamID,
			&i.Price,
			&i.ListedAt,
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

const listTransfersByTeamId = `-- name: ListTransfersByTeamId :many
SELECT id, player_id, seller_team_id, price, listed_at FROM transfers WHERE seller_team_id = $1 AND id > $2 ORDER BY id LIMIT $3
`

type ListTransfersByTeamIdParams struct {
	SellerTeamID int64 `json:"seller_team_id"`
	ID           int64 `json:"id"`
	Limit        int32 `json:"limit"`
}

func (r *pgTransferRepository) ListTransfersByTeamId(
	ctx context.Context,
	arg ListTransfersByTeamIdParams,
) ([]Transfer, error) {
	rows, err := r.db.Query(ctx, listTransfersByTeamId, arg.SellerTeamID, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.PlayerID,
			&i.SellerTeamID,
			&i.Price,
			&i.ListedAt,
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

const updateTransferPriceByIDAndUserID = `-- name: UpdateTransferPriceByIDAndUserID :exec
UPDATE transfers SET price = $2 FROM teams WHERE transfers.id = $1 AND transfers.seller_team_id = teams.id AND teams.user_id = $3
`

type UpdateTransferPriceByIDAndUserIDParams struct {
	ID     int64 `json:"id"`
	Price  int64 `json:"price"`
	UserID int64 `json:"user_id"`
}

func (r *pgTransferRepository) UpdateTransferPriceByIDAndUserID(ctx context.Context, arg UpdateTransferPriceByIDAndUserIDParams) error {
	res, err := r.db.Exec(ctx, updateTransferPriceByIDAndUserID, arg.ID, arg.Price, arg.UserID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *pgTransferRepository) BuyPlayer(
	ctx context.Context,
	transferId int64,
	buyerUserId int64,
) error {
	tx, err := r.db.BeginTx(
		ctx,
		pgx.TxOptions{IsoLevel: pgx.Serializable, AccessMode: pgx.ReadWrite},
	)
	if err != nil {
		return err
	}

	// 0. validation 
	buyerTeam, err := getTeamByUserIDWithQuerier(ctx, tx, buyerUserId)
	if err != nil {
		return postgres.Rollback(ctx, tx, err)
	}

	currentTransfer, err := selectTransferByIdWithQuerier(ctx, tx, transferId)
	// make sure you are not buying from yourself
	if currentTransfer.SellerTeamID == buyerTeam.ID {
		return postgres.Rollback(ctx, tx, ErrConflict)
	}
	// make sure you have enough budget
	if currentTransfer.Price > buyerTeam.Budget {
		return postgres.Rollback(ctx, tx, ErrViolation)
	}

	// 1. delete transfer
	err = deleteTransferByIDWithQuerier(ctx, tx, transferId)
	if err != nil {
		return postgres.Rollback(ctx, tx, err)
	}

	// 2. payment
	// 2.1 charge buyer
	if err := addTeamBudgetWithQuerier(ctx, tx, AddTeamBudgetParams{
		ID:     buyerTeam.ID,
		Budget: -currentTransfer.Price,
	}); err != nil {
		return postgres.Rollback(ctx, tx, err)
	}

	// 2.2 add money to seller
	if err := addTeamBudgetWithQuerier(ctx, tx, AddTeamBudgetParams{
		ID:     currentTransfer.SellerTeamID,
		Budget: currentTransfer.Price,
	}); err != nil {
		return postgres.Rollback(ctx, tx, err)
	}

	// 3. transfer player
	// 3.1 select player
	player, err := getPlayerByIDWithQuerier(ctx, tx, currentTransfer.PlayerID)
	if err != nil {
		return postgres.Rollback(ctx, tx, err)
	}

	// 3.2 calculate random value rise
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	multiplier := (0.1 + rand.Float64()*0.9) + 1
	newPrice := int64(float64(player.Price) * multiplier)

	// 3.3 transfer player to other team
	if err := updatePlayerPriceAndTeamWithQuerrier(ctx, tx, UpdatePlayerPriceAndTeamParams{
		ID:     player.ID,
		Price:  newPrice,
		TeamID: buyerTeam.ID,
	}); err != nil {
		return postgres.Rollback(ctx, tx, err)
	}

	// TODO: insert transfer record
	// 4. insert into transfer_records
	// INSERT INTO transfer_records (id, player_id, seller_team_id, buyer_team_id, sold_price, listed_at)
	// VALUES (
	//     ?,
	//     player_id,
	//     seller_team_id,
	//     ?,
	//     price,
	//     listed_at
	// );

	return tx.Commit(ctx)
}

