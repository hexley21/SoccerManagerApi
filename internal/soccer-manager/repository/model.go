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
	Player struct {
		ID           int64          `json:"id"`
		TeamID       pgtype.Int8    `json:"team_id"`
		CountryCode  pgtype.Text    `json:"country_code"`
		Age          int32          `json:"age"`
		PositionCode string         `json:"position_code"`
		Price        pgtype.Numeric `json:"price"`
	}

	PlayerTranslation struct {
		ID        int64  `json:"id"`
		PlayerID  int64  `json:"player_id"`
		Locale    string `json:"locale"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
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
