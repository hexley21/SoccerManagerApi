package transfer

import (
	"time"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/shopspring/decimal"
)

type transferResponseDTO struct {
	ID           int64     `json:"id"`
	PlayerID     int64     `json:"player_id"`
	SellerTeamID int64     `json:"seller_team_id"`
	Price        int64     `json:"price"`
	ListedAt     time.Time `json:"listed_at"`
} // @name TransferResponse

func transferResponseAdapter(model domain.Transfer) transferResponseDTO {
	return transferResponseDTO{
		ID:           model.ID,
		PlayerID:     model.PlayerID,
		SellerTeamID: model.SellerTeamID,
		Price:        model.Price,
		ListedAt:     model.ListedAt,
	}
}

type createTransferRequestDTO struct {
	PlayerID int64           `json:"player_id" validate:"required"`
	Price    decimal.Decimal `json:"price"     validate:"required,dgte=1"`
} // @name CreateTransferRequest

type updateTransferRequestDTO struct {
	Price decimal.Decimal `json:"price" validate:"required,dgte=1"`
} // @name UpdateTransferRequest
