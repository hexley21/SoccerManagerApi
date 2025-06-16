package repository

import "github.com/jackc/pgx/v5/pgtype"

// Globe
type (
	Country struct {
		Code string
		Name pgtype.Text
	}

	Locale struct {
		Code string
		Name pgtype.Text
	}
)

type (
	User struct {
		ID       int64
		Username string
		Role     string
		Hash     string
	}
)

type (
	Position struct {
		Code        string `json:"code"`
		DefaultName string `json:"default_name"`
	}

	PositionTranslation struct {
		PositionCode string `json:"position_code"`
		Locale       string `json:"locale"`
		Label        string `json:"label"`
	}
)

type (
	Team struct {
		ID           int64          `json:"id"`
		UserID       int64          `json:"user_id"`
		Name         string         `json:"name"`
		CountryCode  pgtype.Text    `json:"country_code"`
		Budget       pgtype.Numeric `json:"budget"`
		TotalPlayers int32          `json:"total_players"`
	}

	TeamTranslation struct {
		TeamID int64  `json:"team_id"`
		Locale string `json:"locale"`
		Name   string `json:"name"`
	}
)
