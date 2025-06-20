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

//go:generate mockgen -destination=mock/mock_team.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service TeamService
type TeamService interface {
	CreateTeam(
		ctx context.Context,
		userId int64,
		name string,
		countryCode domain.CountryCode,
		budget int64,
	) (domain.Team, error)

	GetTeams(
		ctx context.Context,
		locale domain.LocaleCode,
		cursor int64,
		limit int32,
	) ([]domain.Team, error)
	GetTeamById(ctx context.Context, locale domain.LocaleCode, id int64) (domain.Team, error)
	GetTeamByUserId(ctx context.Context, locale domain.LocaleCode, id int64) (domain.Team, error)
	UpdateTeamCountry(
		ctx context.Context,
		userId int64,
		name string,
		countryCode domain.CountryCode,
	) error

	GetTeamTranslations(ctx context.Context, userID int64) (domain.TeamTranslation, error)
	CreateTeamTranslation(
		ctx context.Context,
		locale domain.LocaleCode,
		userID int64,
		name string,
	) error
	UpdateTeamTranslation(
		ctx context.Context,
		locale domain.LocaleCode,
		userID int64,
		name string,
	) error
	DeleteTeamTranslation(ctx context.Context, locale domain.LocaleCode, userID int64) error

	GetAvailableLocales(ctx context.Context, teamID int64) ([]string, error)
}

type teamServiceImpl struct {
	teamRepo            repository.TeamRepository
	teamTranslationRepo repository.TeamTranslationsRepository
}

func NewTeamService(
	teamRepo repository.TeamRepository,
	teamTranslationRepo repository.TeamTranslationsRepository,
) *teamServiceImpl {
	return &teamServiceImpl{
		teamRepo:            teamRepo,
		teamTranslationRepo: teamTranslationRepo,
	}
}

// CreateTeam writes team into database
func (s *teamServiceImpl) CreateTeam(
	ctx context.Context,
	userId int64,
	name string,
	countryCode domain.CountryCode,
	budget int64,
) (domain.Team, error) {
	team, err := s.teamRepo.InsertTeam(ctx, repository.InsertTeamParams{
		UserID:       userId,
		Name:         name,
		CountryCode:  string(countryCode),
		Budget:       budget,
		TotalPlayers: 0,
	})
	if err != nil {
		return domain.Team{}, err
	}

	return domain.TeamAdapter(team), nil
}

// GetTeams returns a list of team by cursor and limit with translation
func (s *teamServiceImpl) GetTeams(
	ctx context.Context,
	locale domain.LocaleCode,
	cursor int64,
	limit int32,
) ([]domain.Team, error) {
	var teams []repository.Team
	var err error

	if locale.Valid() {
		teams, err = s.teamTranslationRepo.ListTranslatedTeamsCursor(
			ctx,
			repository.ListTranslatedTeamsCursorParams{
				ID:     cursor,
				Locale: string(locale),
				Limit:  limit,
			},
		)
	} else {
		teams, err = s.teamRepo.ListTeamsCursor(ctx, repository.ListTeamsCursorParams{
			ID:    cursor,
			Limit: limit,
		})
	}
	if err != nil {
		return nil, err
	}

	res := make([]domain.Team, len(teams))
	for i, t := range teams {
		res[i] = domain.TeamAdapter(t)
	}

	return res, nil
}

// GetTeamById returns a single record with translation
//
// If not found - ErrTeamNotFound
func (s *teamServiceImpl) GetTeamById(
	ctx context.Context,
	locale domain.LocaleCode,
	id int64,
) (domain.Team, error) {
	var team repository.Team
	var err error

	if locale.Valid() {
		team, err = s.teamTranslationRepo.GetTranslatedTeamWithId(
			ctx,
			repository.GetTranslatedTeamWithIdParams{
				ID:     id,
				Locale: string(locale),
			},
		)
	} else {
		team, err = s.teamRepo.GetTeamByID(ctx, id)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, ErrTeamNotFound
		}

		return domain.Team{}, err
	}

	return domain.TeamAdapter(team), nil
}

// GetTeamByUserId finds team by user id with translation
//
// If not found - ErrTeamNotFound
func (s *teamServiceImpl) GetTeamByUserId(
	ctx context.Context,
	locale domain.LocaleCode,
	id int64,
) (domain.Team, error) {
	var team repository.Team
	var err error

	if locale.Valid() {
		team, err = s.teamTranslationRepo.GetTranslatedTeamWithUserId(
			ctx,
			repository.GetTranslatedTeamWithUserIdParams{
				UserID: id,
				Locale: string(locale),
			},
		)
	} else {
		team, err = s.teamRepo.GetTeamByUserID(ctx, id)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, ErrTeamNotFound
		}

		return domain.Team{}, err
	}

	return domain.TeamAdapter(team), nil
}

// UpdateTeamCountry updates country value in db
//
// If not found - ErrTeamNotFound
// If invalid country - ErrNonexistentCode
func (s *teamServiceImpl) UpdateTeamCountry(
	ctx context.Context,
	userId int64,
	name string,
	countryCode domain.CountryCode,
) error {
	if err := s.teamRepo.UpdateTeamDataByUserID(ctx, repository.UpdateTeamDataByUserIDParams{
		UserID:      userId,
		Name:        name,
		CountryCode: string(countryCode),
	}); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTeamNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return ErrNonexistentCode
		}

		return err
	}

	return nil
}

// GetTeamTranslations returns a map of translations
func (s *teamServiceImpl) GetTeamTranslations(
	ctx context.Context,
	userId int64,
) (domain.TeamTranslation, error) {
	translations, err := s.teamTranslationRepo.GetTranslationsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	res := make(domain.TeamTranslation, len(translations))
	for _, tr := range translations {
		res[domain.LocaleCode(tr.Locale)] = tr.Name
	}

	return res, nil
}

// CreateTeamTranslation inserts a translation into db
//
// If invalid locale - ErrNonexistentCode
// If translation exists - ErrTranslationExists
func (s *teamServiceImpl) CreateTeamTranslation(
	ctx context.Context,
	locale domain.LocaleCode,
	userID int64,
	name string,
) error {
	err := s.teamTranslationRepo.InsertTranslationByUserID(
		ctx,
		repository.InsertTranslationByUserIDParams{
			UserID: userID,
			Locale: string(locale),
			Name:   name,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation: // Foreign key violation
				return ErrNonexistentCode
			case pgerrcode.UniqueViolation: // Unique violation
				return ErrTranslationExists
			}
		}
		return err
	}
	return nil
}

// UpdateTeamTranslation updates locale translation label in db
//
// If not found - ErrTranslationNotFound
func (s *teamServiceImpl) UpdateTeamTranslation(
	ctx context.Context,
	locale domain.LocaleCode,
	userID int64,
	name string,
) error {
	err := s.teamTranslationRepo.UpdateTranslationNameByUserID(
		ctx,
		repository.UpdateTranslationNameByUserIDParams{
			UserID: userID,
			Locale: string(locale),
			Name:   name,
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

// DeleteTeamTranslation deletes locale translation from db
//
// If not found - ErrTranslationNotFound
func (s *teamServiceImpl) DeleteTeamTranslation(
	ctx context.Context,
	locale domain.LocaleCode,
	userID int64,
) error {
	err := s.teamTranslationRepo.DeleteTranslationByUserID(
		ctx,
		repository.DeleteTranslationByUserIDParams{
			UserID: userID,
			Locale: string(locale),
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

// GetAvailableLocales returns a list of available locale translations
func (s *teamServiceImpl) GetAvailableLocales(ctx context.Context, teamID int64) ([]string, error) {
	locales, err := s.teamTranslationRepo.ListLocalesByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return locales, nil
}
