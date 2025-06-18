package player

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
)

type playerResponseDTO struct {
	ID           int64                     `json:"id"`
	TeamID       int64                     `json:"team_id"`
	CountryCode  domain.CountryCode        `json:"country_code"`
	FirstName    string                    `json:"first_name"`
	LastName     string                    `json:"last_name"`
	Age          int32                     `json:"age"`
	PositionCode domain.PlayerPositionCode `json:"position_code"`
	Price        int64                     `json:"price"`
} // @name PlayerResponse

func playerResponseAdapter(model domain.Player) playerResponseDTO {
	return playerResponseDTO{
		ID:           model.ID,
		TeamID:       model.TeamID,
		CountryCode:  model.CountryCode,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		Age:          model.Age,
		PositionCode: model.PositionCode,
		Price:        model.Price,
	}
}

type updatePlayerDataRequestDTO struct {
	CountryCode domain.CountryCode `json:"country_code" validate:"required,countrycode"`
	FirstName   string             `json:"first_name"   validate:"required,alphaunicode"`
	LastName    string             `json:"last_name"    validate:"required,alphaunicode"`
} // @name UpdatePlayerDataRequest
