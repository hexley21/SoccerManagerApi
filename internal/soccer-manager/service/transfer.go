package service

import (
	"context"
	"errors"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate mockgen -destination=mock/mock_transfer.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service TransferService
type TransferService interface {
	ListTransfers(
		ctx context.Context,
		id int64,
		limit int32,
	) ([]domain.Transfer, error)
	GetTransferByID(
		ctx context.Context,
		transferID int64,
	) (domain.Transfer, error)
	ListTransfersByTeamId(
		ctx context.Context,
		sellerTeamID, id int64,
		limit int32,
	) ([]domain.Transfer, error)
	GetTransferByPlayerID(
		ctx context.Context,
		playerID int64,
	) (domain.Transfer, error)
	CreateTransfer(
		ctx context.Context,
		userID int64,
		playerID int64,
		price int64,
	) (int64, error)
	DeleteTransfer(ctx context.Context, id int64, userId int64) error
	UpdateTransferPrice(
		ctx context.Context,
		ID int64,
		userId int64,
		price int64,
	) error
	BuyPlayer(
		ctx context.Context,
		transferId int64,
		buyerTeamId int64,
	) error
}

type transferServiceImpl struct {
	transferRepo repository.TransferRepository
}

func NewTransferService(transferRepo repository.TransferRepository) *transferServiceImpl {
	return &transferServiceImpl{
		transferRepo: transferRepo,
	}
}

func (s *transferServiceImpl) ListTransfers(
	ctx context.Context,
	id int64,
	limit int32,
) ([]domain.Transfer, error) {
	transfers, err := s.transferRepo.ListTransfers(
		ctx,
		repository.ListTransfersParams{ID: id, Limit: limit},
	)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Transfer, len(transfers))
	for i, tr := range transfers {
		res[i] = domain.TransferAdapter(tr)
	}

	return res, nil
}

func (s *transferServiceImpl) GetTransferByID(
	ctx context.Context,
	transferID int64,
) (domain.Transfer, error) {
	transfer, err := s.transferRepo.SelectTransferById(ctx, transferID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transfer{}, ErrTransferNotFound
		}
	}
	return domain.TransferAdapter(transfer), nil
}

func (s *transferServiceImpl) ListTransfersByTeamId(
	ctx context.Context,
	sellerTeamID, id int64,
	limit int32,
) ([]domain.Transfer, error) {
	transfers, err := s.transferRepo.ListTransfersByTeamId(
		ctx,
		repository.ListTransfersByTeamIdParams{SellerTeamID: sellerTeamID, ID: id, Limit: limit},
	)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Transfer, len(transfers))
	for i, tr := range transfers {
		res[i] = domain.TransferAdapter(tr)
	}

	return res, nil
}

func (s *transferServiceImpl) GetTransferByPlayerID(
	ctx context.Context,
	playerID int64,
) (domain.Transfer, error) {
	transfer, err := s.transferRepo.GetTransferByPlayerID(ctx, playerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transfer{}, ErrTransferNotFound
		}
	}
	return domain.TransferAdapter(transfer), nil
}

func (s *transferServiceImpl) CreateTransfer(
	ctx context.Context,
	userID int64,
	playerID int64,
	price int64,
) (int64, error) {
	transferId, err := s.transferRepo.InsertTransferRecordByUser(
		ctx,
		repository.InsertTransferRecordByUserParams{
			UserID:   userID,
			PlayerID: playerID,
			Price:    price,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				return 0, ErrNonexistentCode
			case pgerrcode.UniqueViolation:
				return 0, ErrPlayerAlreadyInTransfers
			case pgerrcode.CheckViolation:
				return 0, ErrInvalidArguments
			}
		}

		return 0, err
	}

	return transferId, nil
}

func (s *transferServiceImpl) DeleteTransfer(
	ctx context.Context,
	id int64,
	userId int64,
) error {
	if err := s.transferRepo.DeleteTransferByIDAndUserID(ctx, repository.DeleteTransferByIDAndUserIDParams{
		ID:     id,
		UserID: userId,
	}); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTransferNotFound
		}

		return err
	}

	return nil
}

func (s *transferServiceImpl) UpdateTransferPrice(
	ctx context.Context,
	ID int64,
	userId int64,
	price int64,
) error {
	if err := s.transferRepo.UpdateTransferPriceByIDAndUserID(
		ctx,
		repository.UpdateTransferPriceByIDAndUserIDParams{ID: ID, UserID: userId, Price: price},
	); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTransferNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.CheckViolation:
				return ErrInvalidArguments
			}
		}

		return err
	}

	return nil
}

// BuyPlayer establishes a transaction between seller and buyer
// the player is being transfered to buyer's team
// the transfer record is being inserted
//
// If transfer not found - ErrTransferNotFound
// If buy attempt from yourself - ErrCantBuyFromYourself
// If buy attempt without money - ErrNotEnoughFunds
func (s *transferServiceImpl) BuyPlayer(
	ctx context.Context,
	transferId int64,
	buyerTeamId int64,
) error {
	if err := s.transferRepo.BuyPlayer(ctx, transferId, buyerTeamId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTransferNotFound
		}
		if errors.Is(err, repository.ErrConflict) {
			return ErrCantBuyFromYourself
		}
		if errors.Is(err, repository.ErrViolation) {
			return ErrNotEnoughFunds
		}

		return err
	}

	return nil
}
