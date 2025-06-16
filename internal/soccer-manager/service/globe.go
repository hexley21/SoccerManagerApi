package service

import (
	"context"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
)

type GlobeService interface {
	ListLocales(ctx context.Context) ([]domain.LocaleCode, error)
	ListCountries(ctx context.Context) ([]domain.CountryCode, error)
}

type globeServiceImpl struct {
	globeRepo repository.GlobeRepo
}

func NewGlobeService(globeRepo repository.GlobeRepo) *globeServiceImpl {
	return &globeServiceImpl{globeRepo: globeRepo}
}

func (s *globeServiceImpl) ListLocales(ctx context.Context) ([]domain.LocaleCode, error) {
	locales, err := s.globeRepo.GetAllLocales(ctx)
	if err != nil {
		return []domain.LocaleCode{}, err
	}

	res := make([]domain.LocaleCode, len(locales))

	for i, locs := range locales {
		res[i] = domain.LocaleCode(locs.Code)
	}

	return res, nil
}

func (s *globeServiceImpl) ListCountries(ctx context.Context) ([]domain.CountryCode, error) {
	locales, err := s.globeRepo.GetAllCountries(ctx)
	if err != nil {
		return []domain.CountryCode{}, err
	}

	res := make([]domain.CountryCode, len(locales))

	for i, locs := range locales {
		res[i] = domain.CountryCode(locs.Code)
	}

	return res, nil
}
