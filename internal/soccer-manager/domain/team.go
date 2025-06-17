package domain

import "github.com/shopspring/decimal"

type Team struct {
	ID           int64           `json:"id"`
	UserID       int64           `json:"user_id"`
	Name         string          `json:"name"`
	CountryCode  CountryCode     `json:"country_code"`
	Budget       decimal.Decimal `json:"budget"`
	TotalPlayers int32           `json:"total_players"`
}

type TeamTranslation map[LocaleCode]string // @name TeamTranslation
