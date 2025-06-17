package player

import (
	"strconv"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
)

type playerResponseDTO struct {
	ID           string                    `json:"id"`
	TeamID       string                    `json:"team_id"`
	CountryCode  domain.CountryCode        `json:"country_code"`
	FirstName    string                    `json:"first_name"`
	LastName     string                    `json:"last_name"`
	Age          int32                     `json:"age"`
	PositionCode domain.PlayerPositionCode `json:"position_code"`
	Price        string                    `json:"price"`
} // @name PlayerResponse

func playerResponseAdapter(model domain.Player) playerResponseDTO {
	return playerResponseDTO{
		ID:           strconv.FormatInt(model.ID, 10),
		TeamID:       strconv.FormatInt(model.TeamID, 10),
		CountryCode:  model.CountryCode,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		Age:          model.Age,
		PositionCode: model.PositionCode,
		Price:        model.Price.StringFixed(2),
	}
}

type updatePlayerDataRequestDTO struct {
	CountryCode domain.CountryCode `json:"country_code" validate:"required,countrycode"`
	FirstName   string             `json:"first_name"   validate:"required,alphaunicode"`
	LastName    string             `json:"last_name"    validate:"required,alphaunicode"`
} // @name UpdatePlayerDataRequest
