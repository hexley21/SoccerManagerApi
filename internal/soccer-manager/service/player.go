package service

import (
	"context"
	"errors"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/repository"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=mock/mock_player.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/service PlayerService
type PlayerService interface {
	GetAllPlayers(
		ctx context.Context,
		cursor int64,
		limit int32,
	) ([]domain.Player, error)
	GetPlayerById(
		ctx context.Context,
		playerID int64,
	) (domain.Player, error)
	UpdatePlayerData(
		ctx context.Context,
		userId int64,
		playerID int64,
		firstName string,
		lastName string,
		countryCode domain.CountryCode,
	) error
	GetPlayersByTeamId(
		ctx context.Context,
		teamId int64,
		cursor int64,
		limit int32,
	) ([]domain.Player, error)
	GetPlayersByUserId(
		ctx context.Context,
		userId int64,
		cursor int64,
		limit int32,
	) ([]domain.Player, error)
	CreatePlayer(ctx context.Context, arg CreatePlayerArgs) error
	CreatePlayersBatch(ctx context.Context, args []CreatePlayerArgs) error
}

type playerServiceImpl struct {
	playerRepo repository.PlayerRepository
}

func NewPlayerService(playerRepo repository.PlayerRepository) *playerServiceImpl {
	return &playerServiceImpl{
		playerRepo: playerRepo,
	}
}

func (s *playerServiceImpl) GetAllPlayers(
	ctx context.Context,
	cursor int64,
	limit int32,
) ([]domain.Player, error) {
	players, err := s.playerRepo.ListPlayersByCursor(ctx, repository.ListPlayersByCursorParams{
		ID:    cursor,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	res := make([]domain.Player, len(players))
	for i, p := range players {
		res[i] = domain.PlayerAdapter(p)
	}

	return res, nil
}

// GetPlayerById returns a single player by id
//
// If not found - ErrPlayerNotFound
func (s *playerServiceImpl) GetPlayerById(
	ctx context.Context,
	playerID int64,
) (domain.Player, error) {
	players, err := s.playerRepo.GetPlayerByID(ctx, playerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Player{}, ErrPlayerNotFound
		}

		return domain.Player{}, err
	}

	return domain.PlayerAdapter(players), nil
}

// UpdatePlayerData updates firstname, lastname & countrycode of a team
//
// If not found - ErrPlayerNotFound
func (s *playerServiceImpl) UpdatePlayerData(
	ctx context.Context,
	userId int64,
	playerID int64,
	firstName string,
	lastName string,
	countryCode domain.CountryCode,
) error {
	err := s.playerRepo.UpdatePlayerNameAndCountry(ctx, repository.UpdatePlayerNameAndCountryParams{
		UserID:      userId,
		ID:          playerID,
		FirstName:   firstName,
		LastName:    lastName,
		CountryCode: string(countryCode),
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrPlayerNotFound
		}

		return err
	}

	return nil
}

// GetPlayersByTeamId returns a list of players associated to the team id
func (s *playerServiceImpl) GetPlayersByTeamId(
	ctx context.Context,
	teamId int64,
	cursor int64,
	limit int32,
) ([]domain.Player, error) {
	players, err := s.playerRepo.ListPlayersByTeamID(ctx, repository.ListPlayersByTeamIDParams{
		TeamID: teamId,
		ID:     cursor,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	res := make([]domain.Player, len(players))
	for i, p := range players {
		res[i] = domain.PlayerAdapter(p)
	}

	return res, nil
}

// GetPlayersByUserId returns a list of players associated to the user id
func (s *playerServiceImpl) GetPlayersByUserId(
	ctx context.Context,
	userId int64,
	cursor int64,
	limit int32,
) ([]domain.Player, error) {
	players, err := s.playerRepo.ListPlayersByUserID(ctx, repository.ListPlayersByUserIDParams{
		UserID: userId,
		ID:     cursor,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	res := make([]domain.Player, len(players))
	for i, p := range players {
		res[i] = domain.PlayerAdapter(p)
	}

	return res, nil
}

type CreatePlayerArgs struct {
	TeamID       int64
	CountryCode  domain.CountryCode
	FirstName    string
	LastName     string
	Age          int32
	PositionCode domain.PlayerPositionCode
	Price        int64
}

func NewCreatePlayerArgs(
	teamID int64,
	countryCode domain.CountryCode,
	firstName string,
	lastName string,
	age int32,
	positionCode domain.PlayerPositionCode,
	price int64,
) CreatePlayerArgs {
	return CreatePlayerArgs{
		TeamID:       teamID,
		CountryCode:  countryCode,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		PositionCode: positionCode,
		Price:        price,
	}
}

// CreatePlayersBatch - inserts single player into db
func (s *playerServiceImpl) CreatePlayer(ctx context.Context, arg CreatePlayerArgs) error {
	return s.playerRepo.InsertPlayer(ctx, repository.InsertPlayerParams{
		TeamID:       arg.TeamID,
		CountryCode:  string(arg.CountryCode),
		FirstName:    arg.FirstName,
		LastName:     arg.LastName,
		Age:          arg.Age,
		PositionCode: string(arg.PositionCode),
		Price:        arg.Price,
	})
}

// CreatePlayersBatch - inserts many players at once
func (s *playerServiceImpl) CreatePlayersBatch(ctx context.Context, args []CreatePlayerArgs) error {
	repoParams := make([]repository.InsertPlayerParams, len(args))
	for i, a := range args {
		repoParams[i] = repository.InsertPlayerParams{
			TeamID:       a.TeamID,
			CountryCode:  string(a.CountryCode),
			FirstName:    a.FirstName,
			LastName:     a.LastName,
			Age:          a.Age,
			PositionCode: string(a.PositionCode),
			Price:        a.Price,
		}
	}

	return s.playerRepo.InsertPlayersBatch(ctx, repoParams)
}
