package team

import (
	"strconv"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
)

type teamResponseDTO struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	CountryCode  string `json:"country_code"`
	Budget       string `json:"budget"`
	TotalPlayers int32  `json:"total_players"`
} // @name TeamResponse

func TeamResponseAdapter(team domain.Team) teamResponseDTO {
	return teamResponseDTO{
		ID:           strconv.FormatInt(team.ID, 10),
		UserID:       strconv.FormatInt(team.UserID, 10),
		Name:         team.Name,
		CountryCode:  string(team.CountryCode),
		Budget:       team.Budget.StringFixed(2),
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
