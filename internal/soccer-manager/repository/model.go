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
