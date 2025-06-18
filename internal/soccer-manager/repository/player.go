package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/hexley21/soccer-manager/pkg/infra/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mock/mock_player.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository PlayerRepository
type PlayerRepository interface {
	GetPlayerByID(ctx context.Context, id int64) (Player, error)
	ListPlayersByCursor(ctx context.Context, arg ListPlayersByCursorParams) ([]Player, error)
	ListPlayersByTeamID(ctx context.Context, arg ListPlayersByTeamIDParams) ([]Player, error)
	ListPlayersByUserID(ctx context.Context, arg ListPlayersByUserIDParams) ([]Player, error)
	UpdatePlayerNameAndCountry(ctx context.Context, arg UpdatePlayerNameAndCountryParams) error
	UpdatePlayerPriceAndTeam(ctx context.Context, arg UpdatePlayerPriceAndTeamParams) error
	InsertPlayer(ctx context.Context, arg InsertPlayerParams) error
	InsertPlayersBatch(ctx context.Context, args []InsertPlayerParams) error
}

type pgPlayerRepository struct {
	db            *pgxpool.Pool
	snowflakeNode *snowflake.Node
}

func NewPlayerRepository(db *pgxpool.Pool, snowflakeNode *snowflake.Node) *pgPlayerRepository {
	return &pgPlayerRepository{
		db:            db,
		snowflakeNode: snowflakeNode,
	}
}

const getPlayerByID = `-- name: GetPlayerByID :one
SELECT id, team_id, country_code, first_name, last_name, age, position_code, price FROM players WHERE id = $1
`

func (r *pgPlayerRepository) GetPlayerByID(ctx context.Context, id int64) (Player, error) {
	return getPlayerByIDWithQuerier(ctx, r.db, id)
}

func getPlayerByIDWithQuerier(
	ctx context.Context,
	querier postgres.Querier,
	id int64,
) (Player, error) {
	row := querier.QueryRow(ctx, getPlayerByID, id)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.CountryCode,
		&i.FirstName,
		&i.LastName,
		&i.Age,
		&i.PositionCode,
		&i.Price,
	)
	return i, err
}

const listPlayersByCursor = `-- name: ListPlayersByCursor :many
SELECT id, team_id, country_code, first_name, last_name, age, position_code, price FROM players WHERE id > $1 ORDER BY id LIMIT $2
`

type ListPlayersByCursorParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

func (r *pgPlayerRepository) ListPlayersByCursor(
	ctx context.Context,
	arg ListPlayersByCursorParams,
) ([]Player, error) {
	rows, err := r.db.Query(ctx, listPlayersByCursor, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Player{}
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.CountryCode,
			&i.FirstName,
			&i.LastName,
			&i.Age,
			&i.PositionCode,
			&i.Price,
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

const listPlayersByTeamID = `-- name: ListPlayersByTeamID :many
SELECT p.id, p.team_id, p.country_code, p.first_name, p.last_name, p.age, p.position_code, p.price FROM players p WHERE p.team_id = $1 AND p.id > $2 ORDER BY p.id LIMIT $3
`

type ListPlayersByTeamIDParams struct {
	TeamID int64 `json:"team_id"`
	ID     int64 `json:"id"`
	Limit  int32 `json:"limit"`
}

func (r *pgPlayerRepository) ListPlayersByTeamID(
	ctx context.Context,
	arg ListPlayersByTeamIDParams,
) ([]Player, error) {
	tId := new(pgtype.Int8)
	if err := tId.Scan(arg.TeamID); err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, listPlayersByTeamID, *tId, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Player{}
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.CountryCode,
			&i.FirstName,
			&i.LastName,
			&i.Age,
			&i.PositionCode,
			&i.Price,
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

const listPlayersByUserID = `-- name: ListPlayersByUserID :many
SELECT p.id, p.team_id, p.country_code, p.first_name, p.last_name, p.age, p.position_code, p.price FROM players p JOIN teams t ON p.team_id = t.id WHERE t.user_id = $1 AND p.id > $2 ORDER BY p.id LIMIT $3
`

type ListPlayersByUserIDParams struct {
	UserID int64 `json:"user_id"`
	ID     int64 `json:"id"`
	Limit  int32 `json:"limit"`
}

func (r *pgPlayerRepository) ListPlayersByUserID(
	ctx context.Context,
	arg ListPlayersByUserIDParams,
) ([]Player, error) {
	rows, err := r.db.Query(ctx, listPlayersByUserID, arg.UserID, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Player{}
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.CountryCode,
			&i.FirstName,
			&i.LastName,
			&i.Age,
			&i.PositionCode,
			&i.Price,
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

const updatePlayerNameAndCountry = `-- name: UpdatePlayerNameAndCountry :exec
UPDATE players SET first_name = $3, last_name = $4, country_code = $5 FROM teams WHERE players.team_id = teams.id AND players.id = $2 AND teams.user_id = $1
`
type UpdatePlayerNameAndCountryParams struct {
	UserID      int64  `json:"user_id"`
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CountryCode string `json:"country_code"`
}

func (r *pgPlayerRepository) UpdatePlayerNameAndCountry(
	ctx context.Context,
	arg UpdatePlayerNameAndCountryParams,
) error {
	cc := new(pgtype.Text)
	if err := cc.Scan(arg.CountryCode); err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, updatePlayerNameAndCountry,
		arg.UserID,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		*cc,
	)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const updatePlayerPriceAndTeam = `-- name: UpdatePlayerPriceAndTeam :exec
UPDATE players SET price = $2, team_id = $3 WHERE id = $1
`

type UpdatePlayerPriceAndTeamParams struct {
	ID     int64 `json:"id"`
	Price  int64 `json:"price"`
	TeamID int64 `json:"team_id"`
}

func (r *pgPlayerRepository) UpdatePlayerPriceAndTeam(
	ctx context.Context,
	arg UpdatePlayerPriceAndTeamParams,
) error {
	return updatePlayerPriceAndTeamWithQuerrier(ctx, r.db, arg)
}

func updatePlayerPriceAndTeamWithQuerrier(
	ctx context.Context,
	querier postgres.Querier,
	arg UpdatePlayerPriceAndTeamParams,
) error {
	tId := new(pgtype.Int8)
	if err := tId.Scan(arg.TeamID); err != nil {
		return err
	}

	res, err := querier.Exec(ctx, updatePlayerPriceAndTeam, arg.ID, arg.Price, *tId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const insertPlayer = `-- name: InsertPlayer :exec
INSERT INTO players (id, team_id, country_code, first_name, last_name, age, position_code, price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type InsertPlayerParams struct {
	TeamID       int64  `json:"team_id"`
	CountryCode  string `json:"country_code"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Age          int32  `json:"age"`
	PositionCode string `json:"position_code"`
	Price        int64  `json:"price"`
}

func (r *pgPlayerRepository) insertPlayerWithQuerier(
	ctx context.Context,
	querier postgres.Querier,
	arg InsertPlayerParams,
) error {
	cc := new(pgtype.Text)
	if err := cc.Scan(arg.CountryCode); err != nil {
		return err
	}
	tId := new(pgtype.Int8)
	if err := tId.Scan(arg.TeamID); err != nil {
		return err
	}

	_, err := querier.Exec(ctx, insertPlayer,
		r.snowflakeNode.Generate().Int64(),
		*tId,
		*cc,
		arg.FirstName,
		arg.LastName,
		arg.Age,
		arg.PositionCode,
		arg.Price,
	)
	return err
}

func (r *pgPlayerRepository) InsertPlayer(ctx context.Context, arg InsertPlayerParams) error {
	return r.insertPlayerWithQuerier(ctx, r.db, arg)
}

func (r *pgPlayerRepository) InsertPlayersBatch(
	ctx context.Context,
	args []InsertPlayerParams,
) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	for _, arg := range args {
		if err := r.insertPlayerWithQuerier(ctx, tx, arg); err != nil {
			return postgres.Rollback(ctx, tx, err)
		}
	}

	return tx.Commit(ctx)
}
