package domain

import "github.com/hexley21/soccer-manager/internal/soccer-manager/repository"

type Team struct {
	ID           int64       `json:"id"`
	UserID       int64       `json:"user_id"`
	Name         string      `json:"name"`
	CountryCode  CountryCode `json:"country_code"`
	Budget       int64       `json:"budget"`
	TotalPlayers int32       `json:"total_players"`
}

func TeamAdapter(model repository.Team) Team {
	return Team{
		ID:           model.ID,
		UserID:       model.UserID,
		Name:         model.Name,
		CountryCode:  CountryCode(model.CountryCode.String),
		Budget:       model.Budget,
		TotalPlayers: model.TotalPlayers,
	}
}

type TeamTranslation map[LocaleCode]string // @name TeamTranslation
