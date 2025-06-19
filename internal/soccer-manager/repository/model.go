package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

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
	Team struct {
		ID           int64
		UserID       int64
		Name         string
		CountryCode  pgtype.Text
		Budget       int64
		TotalPlayers int32
	}

	TeamTranslation struct {
		TeamID int64
		Locale string
		Name   string
	}
)

type (
	Position struct {
		Code        string
		DefaultName string
	}

	PositionTranslation struct {
		PositionCode string
		Locale       string
		Label        string
	}
)

type (
	Player struct {
		ID           int64
		TeamID       pgtype.Int8
		CountryCode  pgtype.Text
		FirstName    string
		LastName     string
		Age          int32
		PositionCode string
		Price        int64
	}
)

type (
	Transfer struct {
		ID           int64
		PlayerID     int64
		SellerTeamID int64
		Price        int64
		ListedAt     pgtype.Timestamptz
	}
)

type (
	TransferRecord struct {
		ID           int64
		PlayerID     int64
		SellerTeamID int64
		BuyerTeamID  int64
		SoldPrice    int64
		ListedAt     pgtype.Timestamptz
		SoldAt       pgtype.Timestamptz
	}
)
