package transfer_record

import (
	"time"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
)

type transferRecordResponseDTO struct {
	ID           int64     `json:"id"`
	PlayerID     int64     `json:"player_id"`
	SellerTeamID int64     `json:"seller_team_id"`
	BuyerTeamID  int64     `json:"buyer_team_id"`
	SoldPrice    int64     `json:"sold_price"`
	ListedAt     time.Time `json:"listed_at"`
	SoldAt       time.Time `json:"sold_at"`
} // @name TransferRecordResponse

func transferRecordResponseAdapter(model domain.TransferRecord) transferRecordResponseDTO {
	return transferRecordResponseDTO{
		ID:           model.ID,
		PlayerID:     model.PlayerID,
		SellerTeamID: model.SellerTeamID,
		BuyerTeamID:  model.BuyerTeamID,
		SoldPrice:    model.SoldPrice,
		ListedAt:     model.ListedAt,
		SoldAt:       model.SoldAt,
	}
}
