package service

import (
	"context"
	"errors"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=mock/mock_transfer_record.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service TransferRecordService
type TransferRecordService interface {
	ListTransferRecords(ctx context.Context, id int64, limit int32) ([]domain.TransferRecord, error)
	GetTransferRecordByID(ctx context.Context, id int64) (domain.TransferRecord, error)
}

type transferRecordServiceImpl struct {
	transferRecordRepo repository.TransferRecordRepository
}

func NewTransferRecordService(
	transferRecordRepo repository.TransferRecordRepository,
) *transferRecordServiceImpl {
	return &transferRecordServiceImpl{
		transferRecordRepo: transferRecordRepo,
	}
}

func (s *transferRecordServiceImpl) ListTransferRecords(
	ctx context.Context,
	id int64,
	limit int32,
) ([]domain.TransferRecord, error) {
	records, err := s.transferRecordRepo.ListTransferRecords(
		ctx,
		repository.ListTransferRecordsParams{
			ID:    id,
			Limit: limit,
		},
	)
	if err != nil {
		return nil, err
	}

	res := make([]domain.TransferRecord, len(records))
	for i, tr := range records {
		res[i] = domain.TransferRecordAdapter(tr)
	}

	return res, nil
}

func (s *transferRecordServiceImpl) GetTransferRecordByID(
	ctx context.Context,
	id int64,
) (domain.TransferRecord, error) {
	record, err := s.transferRecordRepo.GetTransferRecordByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TransferRecord{}, ErrTransferRecordNotFound
		}
	}

	return domain.TransferRecordAdapter(record), nil
}
