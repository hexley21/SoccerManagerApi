package domain

type Team struct {
	ID           int64       `json:"id"`
	UserID       int64       `json:"user_id"`
	Name         string      `json:"name"`
	CountryCode  CountryCode `json:"country_code"`
	Budget       int64       `json:"budget"`
	TotalPlayers int32       `json:"total_players"`
}

type TeamTranslation map[LocaleCode]string // @name TeamTranslation
