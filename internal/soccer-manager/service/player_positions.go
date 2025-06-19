package service

import (
	"context"
	"errors"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate mockgen -destination=mock/mock_player_position.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service PlayerPositionService
type PlayerPositionService interface {
	ListPositionCodes(ctx context.Context) ([]domain.PlayerPositionCode, error)
	ListPositions(ctx context.Context) ([]domain.PlayerPosition, error)
	ListPositionTranslationsByLocale(
		ctx context.Context,
		locale domain.LocaleCode,
	) ([]domain.PlayerPosition, error)
	ListPositionTranslationsByCode(
		ctx context.Context,
		code domain.PlayerPositionCode,
	) ([]domain.PlayerPosition, error)
	ListPositionTranslations(ctx context.Context) ([]domain.PlayerPositionWithLocale, error)
	CreateTranslation(
		ctx context.Context,
		code domain.PlayerPositionCode,
		locale domain.LocaleCode,
		label string,
	) error
	UpdateTranslation(
		ctx context.Context,
		code domain.PlayerPositionCode,
		locale domain.LocaleCode,
		label string,
	) error
	DeleteTranslation(
		ctx context.Context,
		code domain.PlayerPositionCode,
		locale domain.LocaleCode,
	) error
}

type playerPositionServiceImpl struct {
	playerPositionRepo repository.PlayerPositionRepository
}

func NewPlayerPositionService(
	playerPositionRepo repository.PlayerPositionRepository,
) *playerPositionServiceImpl {
	return &playerPositionServiceImpl{
		playerPositionRepo: playerPositionRepo,
	}
}

// ListPositionCodes returns a list of position codes
func (s *playerPositionServiceImpl) ListPositionCodes(
	ctx context.Context,
) ([]domain.PlayerPositionCode, error) {
	codes, err := s.playerPositionRepo.GetAllPositionCodes(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]domain.PlayerPositionCode, len(codes))

	for i, code := range codes {
		res[i] = domain.PlayerPositionCode(code)
	}

	return res, nil
}

// ListPositionCodes returns a list of position codes
func (s *playerPositionServiceImpl) ListPositions(
	ctx context.Context,
) ([]domain.PlayerPosition, error) {
	positions, err := s.playerPositionRepo.GetAllPositions(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]domain.PlayerPosition, len(positions))

	for i, pos := range positions {
		res[i] = domain.PlayerPosition{
			Code:  domain.PlayerPositionCode(pos.Code),
			Label: pos.DefaultName,
		}
	}

	return res, nil
}

// ListPositionTranslationsByLocale returns a list of translated position by locale
func (s *playerPositionServiceImpl) ListPositionTranslationsByLocale(
	ctx context.Context,
	locale domain.LocaleCode,
) ([]domain.PlayerPosition, error) {
	translated, err := s.playerPositionRepo.GetPositionTranslationsByLocale(
		ctx,
		string(locale),
	)
	if err != nil {
		return nil, err
	}

	res := make([]domain.PlayerPosition, len(translated))

	for i, trns := range translated {
		res[i] = domain.PlayerPosition{
			Code:  domain.PlayerPositionCode(trns.PositionCode),
			Label: trns.Label,
		}
	}

	return res, nil
}

// ListPositionTranslations returns a list of all available translations
func (s *playerPositionServiceImpl) ListPositionTranslations(
	ctx context.Context,
) ([]domain.PlayerPositionWithLocale, error) {
	translated, err := s.playerPositionRepo.GetPositionTranslations(ctx)
	if err != nil {
		return []domain.PlayerPositionWithLocale{}, err
	}

	res := make([]domain.PlayerPositionWithLocale, len(translated))
	for i, tr := range translated {
		res[i] = domain.PlayerPositionWithLocale{
			PlayerPosition: domain.PlayerPosition{
				Code:  domain.PlayerPositionCode(tr.PositionCode),
				Label: tr.Label,
			},
			Locale: domain.LocaleCode(tr.Locale),
		}
	}

	return res, nil
}

// ListPositionTranslationsByCode returns a list of translations by position code
func (s *playerPositionServiceImpl) ListPositionTranslationsByCode(
	ctx context.Context,
	code domain.PlayerPositionCode,
) ([]domain.PlayerPosition, error) {
	translated, err := s.playerPositionRepo.GetPositionTranslationsByPositionCode(
		ctx,
		string(code),
	)
	if err != nil {
		return nil, err
	}

	res := make([]domain.PlayerPosition, len(translated))

	for i, trns := range translated {
		res[i] = domain.PlayerPosition{
			Code:  domain.PlayerPositionCode(trns.PositionCode),
			Label: trns.Label,
		}
	}

	return res, nil
}

// CreateTranslation creates a translation for provided position code
//
// If invalid locale - ErrNonexistentCode
// If translation exists - ErrTranslationExists
func (s *playerPositionServiceImpl) CreateTranslation(
	ctx context.Context,
	code domain.PlayerPositionCode,
	locale domain.LocaleCode,
	label string,
) error {
	err := s.playerPositionRepo.InsertPositionTranslation(
		ctx,
		repository.InsertPositionTranslationParams{
			PositionCode: string(code),
			Locale:       string(locale),
			Label:        label,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				return ErrNonexistentCode
			case pgerrcode.UniqueViolation:
				return ErrTranslationExists
			}
		}
		return err
	}
	
	return nil
}

// UpdateTranslation changes label fo provided position translation
//
// If translation not found - ErrTranslationNotFound
func (s *playerPositionServiceImpl) UpdateTranslation(
	ctx context.Context,
	code domain.PlayerPositionCode,
	locale domain.LocaleCode,
	label string,
) error {
	err := s.playerPositionRepo.UpdatePositionTranslationLabel(
		ctx,
		repository.UpdatePositionTranslationLabelParams{
			PositionCode: string(code),
			Locale:       string(locale),
			Label:        label,
		},
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTranslationNotFound
		}

		return err
	}

	return nil
}

// DeleteTranslation removes translation by code and locale
//
// If translation not found - ErrTranslationNotFound
func (s *playerPositionServiceImpl) DeleteTranslation(
	ctx context.Context,
	code domain.PlayerPositionCode,
	locale domain.LocaleCode,
) error {
	err := s.playerPositionRepo.DeletePositionTranslation(
		ctx,
		repository.DeletePositionTranslationParams{
			PositionCode: string(code),
			Locale:       string(locale),
		},
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTranslationNotFound
		}

		return err
	}

	return nil
}
