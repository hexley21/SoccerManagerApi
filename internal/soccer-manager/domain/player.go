package domain

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/shopspring/decimal"
)

type Player struct {
	ID           int64              `json:"id"`
	TeamID       int64              `json:"team_id"`
	CountryCode  CountryCode        `json:"country_code"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Age          int32              `json:"age"`
	PositionCode PlayerPositionCode `json:"position_code"`
	Price        decimal.Decimal    `json:"price"`
}

func PlayerAdapter(model repository.Player) Player {
	return Player{
		ID:           model.ID,
		TeamID:       model.TeamID.Int64,
		CountryCode:  CountryCode(model.CountryCode.String),
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		Age:          model.Age,
		PositionCode: PlayerPositionCode(model.PositionCode),
		Price:        decimal.NewFromBigInt(model.Price.Int, model.Price.Exp),
	}
}
