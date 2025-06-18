package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=mock/mock_translations.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/repository TeamTranslationsRepository
type TeamTranslationsRepository interface {
	InsertTranslationByUserID(ctx context.Context, arg InsertTranslationByUserIDParams) error
	GetTranslationsByUserID(ctx context.Context, userID int64) ([]TeamTranslation, error)
	DeleteTranslationByUserID(ctx context.Context, arg DeleteTranslationByUserIDParams) error
	UpdateTranslationNameByUserID(
		ctx context.Context,
		arg UpdateTranslationNameByUserIDParams,
	) error

	// translated team
	ListTranslatedTeamsCursor(
		ctx context.Context,
		arg ListTranslatedTeamsCursorParams,
	) ([]Team, error)
	GetTranslatedTeamWithId(ctx context.Context, arg GetTranslatedTeamWithIdParams) (Team, error)
	GetTranslatedTeamWithUserId(
		ctx context.Context,
		arg GetTranslatedTeamWithUserIdParams,
	) (Team, error)

	ListLocalesByTeamID(ctx context.Context, teamID int64) ([]string, error)
}

type pgTeamTranslations struct {
	db *pgxpool.Pool
}

func NewTeamTranslationsRepository(db *pgxpool.Pool) *pgTeamTranslations {
	return &pgTeamTranslations{
		db: db,
	}
}

const insertTranslationByUserID = `-- name: InsertTranslationByUserID :exec
INSERT INTO team_translations (team_id, locale, name) SELECT t.id, $2, $3 FROM teams t WHERE t.user_id = $1
`

type InsertTranslationByUserIDParams struct {
	UserID int64  `json:"user_id"`
	Locale string `json:"locale"`
	Name   string `json:"name"`
}

func (r *pgTeamTranslations) InsertTranslationByUserID(
	ctx context.Context,
	arg InsertTranslationByUserIDParams,
) error {
	_, err := r.db.Exec(ctx, insertTranslationByUserID, arg.UserID, arg.Locale, arg.Name)
	return err
}

const getTranslationsByUserID = `-- name: GetTranslationsByUserID :many
SELECT tt.team_id, tt.locale, tt.name FROM team_translations tt JOIN teams t ON t.id = tt.team_id WHERE t.user_id = $1
`

func (r *pgTeamTranslations) GetTranslationsByUserID(
	ctx context.Context,
	userID int64,
) ([]TeamTranslation, error) {
	rows, err := r.db.Query(ctx, getTranslationsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TeamTranslation{}
	for rows.Next() {
		var i TeamTranslation
		if err := rows.Scan(&i.TeamID, &i.Locale, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteTranslationByUserID = `-- name: DeleteTranslationByUserID :exec
DELETE FROM team_translations AS tt USING teams AS t WHERE tt.team_id = t.id AND t.user_id = $1 AND tt.locale = $2
`

type DeleteTranslationByUserIDParams struct {
	UserID int64  `json:"user_id"`
	Locale string `json:"locale"`
}

func (r *pgTeamTranslations) DeleteTranslationByUserID(
	ctx context.Context,
	arg DeleteTranslationByUserIDParams,
) error {
	res, err := r.db.Exec(ctx, deleteTranslationByUserID, arg.UserID, arg.Locale)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return err
}

const updateTranslationNameByUserID = `-- name: UpdateTranslationNameByUserID :exec
UPDATE team_translations AS tt SET name = $3 FROM teams AS t WHERE tt.team_id = t.id AND t.user_id = $1 AND tt.locale = $2
`

type UpdateTranslationNameByUserIDParams struct {
	UserID int64  `json:"user_id"`
	Locale string `json:"locale"`
	Name   string `json:"name"`
}

func (r *pgTeamTranslations) UpdateTranslationNameByUserID(
	ctx context.Context,
	arg UpdateTranslationNameByUserIDParams,
) error {
	res, err := r.db.Exec(ctx, updateTranslationNameByUserID, arg.UserID, arg.Locale, arg.Name)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return err
}

const getTranslatedTeamWithId = `-- name: GetTranslatedTeamWithId :one
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
WHERE t.id = $1
`

type GetTranslatedTeamWithIdParams struct {
	ID     int64  `json:"id"`
	Locale string `json:"locale"`
}

func (r *pgTeamTranslations) GetTranslatedTeamWithId(
	ctx context.Context,
	arg GetTranslatedTeamWithIdParams,
) (Team, error) {
	row := r.db.QueryRow(ctx, getTranslatedTeamWithId, arg.ID, arg.Locale)
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

const getTranslatedTeamWithUserId = `-- name: GetTranslatedTeamWithUserId :one
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
WHERE t.user_id = $1
`

type GetTranslatedTeamWithUserIdParams struct {
	UserID int64  `json:"user_id"`
	Locale string `json:"locale"`
}

func (r *pgTeamTranslations) GetTranslatedTeamWithUserId(
	ctx context.Context,
	arg GetTranslatedTeamWithUserIdParams,
) (Team, error) {
	row := r.db.QueryRow(ctx, getTranslatedTeamWithUserId, arg.UserID, arg.Locale)
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

const listTranslatedTeamsCursor = `-- name: ListTranslatedTeamsCursor :many
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
LIMIT $3
`

type ListTranslatedTeamsCursorParams struct {
	ID     int64  `json:"id"`
	Locale string `json:"locale"`
	Limit  int32  `json:"limit"`
}

func (r *pgTeamTranslations) ListTranslatedTeamsCursor(
	ctx context.Context,
	arg ListTranslatedTeamsCursorParams,
) ([]Team, error) {
	rows, err := r.db.Query(ctx, listTranslatedTeamsCursor, arg.ID, arg.Locale, arg.Limit)
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

const listLocalesByTeamID = `-- name: ListLocalesByTeamID :many
SELECT locale FROM team_translations WHERE team_id = $1
`

func (r *pgTeamTranslations) ListLocalesByTeamID(
	ctx context.Context,
	teamID int64,
) ([]string, error) {
	rows, err := r.db.Query(ctx, listLocalesByTeamID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var locale string
		if err := rows.Scan(&locale); err != nil {
			return nil, err
		}
		items = append(items, locale)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
