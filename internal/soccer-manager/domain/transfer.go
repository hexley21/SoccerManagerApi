package domain

import (
	"time"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
)

type Transfer struct {
	ID           int64
	PlayerID     int64
	SellerTeamID int64
	Price        int64
	ListedAt     time.Time
}

func TransferAdapter(model repository.Transfer) Transfer {
	return Transfer{
		ID:           model.ID,
		PlayerID:     model.PlayerID,
		SellerTeamID: model.SellerTeamID,
		Price:        model.Price,
		ListedAt:     model.ListedAt.Time,
	}
}
