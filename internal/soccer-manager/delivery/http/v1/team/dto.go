package team

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
)

type teamResponseDTO struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	CountryCode  string `json:"country_code"`
	Budget       int64  `json:"budget"`
	TotalPlayers int32  `json:"total_players"`
} // @name TeamResponse

func TeamResponseAdapter(team domain.Team) teamResponseDTO {
	return teamResponseDTO{
		ID:           team.ID,
		UserID:       team.UserID,
		Name:         team.Name,
		CountryCode:  string(team.CountryCode),
		Budget:       team.Budget,
		TotalPlayers: team.TotalPlayers,
	}
}

type updateTeamRequestDTO struct {
	Name        string             `json:"name"         validate:"required"`
	CountryCode domain.CountryCode `json:"country_code" validate:"required,countrycode"`
} // @name UpdateTeamRequest

type createTeamTranslationRequestDTO struct {
	Locale domain.LocaleCode `json:"locale" validate:"required,localecode"`
	Name   string            `json:"label"  validate:"required,alphaunicode"`
} // @name CreateTeamTranslationRequest

type updateTeamTranslationRequestDTO struct {
	Locale domain.LocaleCode `json:"locale" validate:"required,localecode"`
	Name   string            `json:"label"  validate:"required,alphaunicode"`
} // @name UpdateTeamTranslationRequest

type deleteTeamTranslationRequestDTO struct {
	Locale domain.LocaleCode `json:"locale" validate:"required,localecode"`
} // @name DeleteTeamTranslationRequest
