package domain

import (
	"time"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
)

type TransferRecord struct {
	ID           int64
	PlayerID     int64
	SellerTeamID int64
	BuyerTeamID  int64
	SoldPrice    int64
	ListedAt     time.Time
	SoldAt       time.Time
}

func TransferRecordAdapter(model repository.TransferRecord) TransferRecord {
	return TransferRecord{
		ID:           model.ID,
		PlayerID:     model.PlayerID,
		SellerTeamID: model.SellerTeamID,
		BuyerTeamID:  model.BuyerTeamID,
		SoldPrice:    model.SoldPrice,
		ListedAt:     model.ListedAt.Time,
		SoldAt:       model.SoldAt.Time,
	}
}
