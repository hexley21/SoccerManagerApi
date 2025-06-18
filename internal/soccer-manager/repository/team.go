package repository

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/hexley21/soccer-manager/pkg/infra/postgres"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mock/mock_team.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository TeamRepository
type TeamRepository interface {
	DeleteTeamByID(ctx context.Context, id int64) error
	GetTeamByID(ctx context.Context, id int64) (Team, error)
	GetTeamByUserID(ctx context.Context, userID int64) (Team, error)
	InsertTeam(ctx context.Context, arg InsertTeamParams) (Team, error)
	ListTeamsCursor(ctx context.Context, arg ListTeamsCursorParams) ([]Team, error)
	UpdateTeamDataByUserID(ctx context.Context, arg UpdateTeamDataByUserIDParams) error
}

type pgTeamRepository struct {
	db            *pgxpool.Pool
	snowflakeNode *snowflake.Node
}

func NewTeamRepository(db *pgxpool.Pool, snowflakeNode *snowflake.Node) *pgTeamRepository {
	return &pgTeamRepository{
		db:            db,
		snowflakeNode: snowflakeNode,
	}
}

const deleteTeamByID = `-- name: DeleteTeamByID :exec
DELETE FROM teams WHERE id = $1
`

func (r *pgTeamRepository) DeleteTeamByID(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, deleteTeamByID, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const getTeamByID = `-- name: GetTeamByID :one
SELECT id, user_id, name, country_code, budget, total_players FROM teams WHERE id = $1
`

func (r *pgTeamRepository) GetTeamByID(ctx context.Context, id int64) (Team, error) {
	row := r.db.QueryRow(ctx, getTeamByID, id)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CountryCode,
		&i.Budget,
		&i.TotalPlayers,
	)
	return i, err
}

const getTeamByUserID = `-- name: GetTeamByUserID :one
SELECT id, user_id, name, country_code, budget, total_players FROM teams WHERE user_id = $1
`

func (r *pgTeamRepository) GetTeamByUserID(ctx context.Context, userID int64) (Team, error) {
	return getTeamByUserIDWithQuerier(ctx, r.db, userID)
}

func getTeamByUserIDWithQuerier(
	ctx context.Context,
	querier postgres.Querier,
	userID int64,
) (Team, error) {
	row := querier.QueryRow(ctx, getTeamByUserID, userID)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CountryCode,
		&i.Budget,
		&i.TotalPlayers,
	)
	return i, err
}

const insertTeam = `-- name: InsertTeam :one
INSERT INTO teams (id, user_id, name, country_code, budget, total_players) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, user_id, name, country_code, budget, total_players
`

type InsertTeamParams struct {
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	CountryCode  string `json:"country_code"`
	Budget       int64  `json:"budget"`
	TotalPlayers int32  `json:"total_players"`
}

func (r *pgTeamRepository) InsertTeam(ctx context.Context, arg InsertTeamParams) (Team, error) {
	cc := new(pgtype.Text)
	if err := cc.Scan(arg.CountryCode); err != nil {
		return Team{}, err
	}

	row := r.db.QueryRow(ctx, insertTeam,
		r.snowflakeNode.Generate().Int64(),
		arg.UserID,
		arg.Name,
		*cc,
		arg.Budget,
		arg.TotalPlayers,
	)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CountryCode,
		&i.Budget,
		&i.TotalPlayers,
	)
	return i, err
}

const listTeamsCursor = `-- name: ListTeamsCursor :many
SELECT id, user_id, name, country_code, budget, total_players FROM teams WHERE id > $1 ORDER BY id LIMIT $2
`

type ListTeamsCursorParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

func (r *pgTeamRepository) ListTeamsCursor(
	ctx context.Context,
	arg ListTeamsCursorParams,
) ([]Team, error) {
	rows, err := r.db.Query(ctx, listTeamsCursor, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Team{}
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.CountryCode,
			&i.Budget,
			&i.TotalPlayers,
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

const updateTeamDataByUserID = `-- name: UpdateTeamDataByUserID :exec
UPDATE teams SET name = $2, country_code = $3 WHERE user_id = $1
`

type UpdateTeamDataByUserIDParams struct {
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
}

func (r *pgTeamRepository) UpdateTeamDataByUserID(
	ctx context.Context,
	arg UpdateTeamDataByUserIDParams,
) error {
	cc := new(pgtype.Text)
	if err := cc.Scan(arg.CountryCode); err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, updateTeamDataByUserID, arg.UserID, arg.Name, *cc)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const addTeamBudget = `-- name: AddTeamBudget :exec
UPDATE teams SET budget = budget + $2 WHERE id = $1
`

type AddTeamBudgetParams struct {
	ID     int64 `json:"id"`
	Budget int64 `json:"budget"`
}

func addTeamBudgetWithQuerier(
	ctx context.Context,
	querier postgres.Querier,
	arg AddTeamBudgetParams,
) error {
	res, err := querier.Exec(ctx, addTeamBudget, arg.ID, arg.Budget)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return err
}
