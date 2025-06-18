package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mock/mock_player_position.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository PlayerPositionRepository
type PlayerPositionRepository interface {
	GetAllPositionCodes(ctx context.Context) ([]string, error)
	GetAllPositions(ctx context.Context) ([]Position, error)
	GetPositionTranslationsByLocale(
		ctx context.Context,
		locale string,
	) ([]PositionTranslation, error)
	GetPositionTranslationsByPositionCode(
		ctx context.Context,
		positionCode string,
	) ([]PositionTranslation, error)
	GetPositionTranslations(ctx context.Context) ([]PositionTranslation, error)
	InsertPositionTranslation(ctx context.Context, arg InsertPositionTranslationParams) error
	UpdatePositionTranslationLabel(
		ctx context.Context,
		arg UpdatePositionTranslationLabelParams,
	) error
	DeletePositionTranslation(
		ctx context.Context,
		arg DeletePositionTranslationParams,
	) error
}

// TODO: add position caching
type pgPlayerPositionRepo struct {
	db *pgxpool.Pool
}

func NewPlayerPositionRepo(
	db *pgxpool.Pool,
) *pgPlayerPositionRepo {
	return &pgPlayerPositionRepo{
		db: db,
	}
}

const getAllPositionCodes = `-- name: GetAllPositionCodes :many
SELECT code FROM positions ORDER BY code
`

func (r *pgPlayerPositionRepo) GetAllPositionCodes(ctx context.Context) ([]string, error) {
	rows, err := r.db.Query(ctx, getAllPositionCodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		items = append(items, code)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllPositions = `-- name: GetAllPositions :many
SELECT code, default_name FROM positions ORDER BY code
`

func (r *pgPlayerPositionRepo) GetAllPositions(ctx context.Context) ([]Position, error) {
	rows, err := r.db.Query(ctx, getAllPositions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Position{}
	for rows.Next() {
		var i Position
		if err := rows.Scan(&i.Code, &i.DefaultName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPositionTranslationsByLocale = `-- name: GetPositionTranslationsByLocale :many
SELECT position_code, locale, label FROM position_translations WHERE locale = $1 ORDER BY locale
`

func (r *pgPlayerPositionRepo) GetPositionTranslationsByLocale(
	ctx context.Context,
	locale string,
) ([]PositionTranslation, error) {
	rows, err := r.db.Query(ctx, getPositionTranslationsByLocale, locale)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PositionTranslation{}
	for rows.Next() {
		var i PositionTranslation
		if err := rows.Scan(&i.PositionCode, &i.Locale, &i.Label); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPositionTranslationsByPositionCode = `-- name: GetPositionTranslationsByPositionCode :many
SELECT position_code, locale, label FROM position_translations WHERE position_code = $1 ORDER BY locale
`

func (r *pgPlayerPositionRepo) GetPositionTranslationsByPositionCode(
	ctx context.Context,
	positionCode string,
) ([]PositionTranslation, error) {
	rows, err := r.db.Query(ctx, getPositionTranslationsByPositionCode, positionCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PositionTranslation{}
	for rows.Next() {
		var i PositionTranslation
		if err := rows.Scan(&i.PositionCode, &i.Locale, &i.Label); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPositionTranslations = `-- name: GetPositionTranslations :many
SELECT position_code, locale, label FROM position_translations ORDER BY position_code
`

func (r *pgPlayerPositionRepo) GetPositionTranslations(
	ctx context.Context,
) ([]PositionTranslation, error) {
	rows, err := r.db.Query(ctx, getPositionTranslations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PositionTranslation{}
	for rows.Next() {
		var i PositionTranslation
		if err := rows.Scan(&i.PositionCode, &i.Locale, &i.Label); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertPositionTranslation = `-- name: InsertPositionTranslation :exec
INSERT INTO position_translations (position_code, locale, label) VALUES ($1, $2, $3)
`

type InsertPositionTranslationParams struct {
	PositionCode string `json:"position_code"`
	Locale       string `json:"locale"`
	Label        string `json:"label"`
}

func (r *pgPlayerPositionRepo) InsertPositionTranslation(
	ctx context.Context,
	arg InsertPositionTranslationParams,
) error {
	_, err := r.db.Exec(ctx, insertPositionTranslation, arg.PositionCode, arg.Locale, arg.Label)
	return err
}

const updatePositionTranslationLabel = `-- name: UpdatePositionTranslationLabel :exec
UPDATE position_translations SET label = $3 WHERE position_code = $1 AND locale = $2
`

type UpdatePositionTranslationLabelParams struct {
	PositionCode string `json:"position_code"`
	Locale       string `json:"locale"`
	Label        string `json:"label"`
}

func (r *pgPlayerPositionRepo) UpdatePositionTranslationLabel(
	ctx context.Context,
	arg UpdatePositionTranslationLabelParams,
) error {
	res, err := r.db.Exec(
		ctx,
		updatePositionTranslationLabel,
		arg.PositionCode,
		arg.Locale,
		arg.Label,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

const deletePositionTranslation = `-- name: DeletePositionTranslation :exec
DELETE FROM position_translations WHERE position_code = $1 AND locale = $2
`

type DeletePositionTranslationParams struct {
	PositionCode string `json:"position_code"`
	Locale       string `json:"locale"`
}

func (r *pgPlayerPositionRepo) DeletePositionTranslation(
	ctx context.Context,
	arg DeletePositionTranslationParams,
) error {
	res, err := r.db.Exec(ctx, deletePositionTranslation, arg.PositionCode, arg.Locale)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
