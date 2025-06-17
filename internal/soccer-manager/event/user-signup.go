package event

import (
	"context"
	"math/rand"
	"time"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/hexley21/soccer-manager/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type userSignUpHandlerImpl struct {
	teamService   service.TeamService
	playerService service.PlayerService
	cfg           config.UserSignUp

	logger echo.Logger

	sem chan struct{}
}

func newUserSignUpHandler(
	teamService service.TeamService,
	playerService service.PlayerService,
	cfg config.UserSignUp,
	logger echo.Logger,
) *userSignUpHandlerImpl {
	return &userSignUpHandlerImpl{
		teamService: teamService,
		playerService: playerService,
		cfg:         cfg,
		logger:      logger,

		sem: make(chan struct{}, cfg.GoroutineCount),
	}
}

func (h *userSignUpHandlerImpl) Acquire() {
	h.sem <- struct{}{}
}

func (h *userSignUpHandlerImpl) Release() {
	<-h.sem
}

// TODO: add random or ip based team assignation
// TODO: add process failure handling
func (h *userSignUpHandlerImpl) Handle(userId int64, username string) {
	h.logger.Debugf("received signal for: %v %v", userId, username)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Errorf("panic recovered while handling user signup: %v", r)
			}
		}()

		h.Acquire()
		defer h.Release()

		ctx, cancel := context.WithTimeout(context.Background(), h.cfg.Timeout)
		defer cancel()

		team, err := h.Createteam(ctx, userId, username, h.cfg.TeamBudget)
		if err != nil {
			h.logger.Errorf("failed to create teams at signup: %v", err)
			return
		}

		// prepare argments for palyer insertion
		args := buildBatchPlayers(
			team.ID,
			team.CountryCode,
			h.cfg.PlayerMinAge,
			h.cfg.PlayerMaxAge,
			h.cfg.PlayerBudget,
			h.cfg.TeamMembers,
		)

		if err := h.playerService.CreatePlayersBatch(ctx, args); err != nil {
			h.logger.Errorf("failed to create players at signup: %v", err)
			return
		}

		h.logger.Debugf("on signup done, team: %v", team)
	}()
}

func (h *userSignUpHandlerImpl) Createteam(
	ctx context.Context,
	userId int64,
	username string,
	budget int64,
) (domain.Team, error) {
	team, err := h.teamService.CreateTeam(
		ctx,
		userId,
		username,
		domain.CountryCode("GE"),
		budget,
	)
	if err != nil {

		return domain.Team{}, err
	}

	return team, nil
}

func buildBatchPlayers(
	teamID int64,
	countryCode domain.CountryCode,
	minAge int,
	maxAge int,
	price int64,
	memberConfig config.TeamMembers,
) []service.CreatePlayerArgs {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	members := map[domain.PlayerPositionCode]int{
		domain.PlayerPositionCodeGLK: memberConfig.GLK,
		domain.PlayerPositionCodeDEF: memberConfig.DEF,
		domain.PlayerPositionCodeMID: memberConfig.MID,
		domain.PlayerPositionCodeATK: memberConfig.ATK,
	}

	var players []service.CreatePlayerArgs
	for position, amount := range members {
		for range amount {
			player := service.NewCreatePlayerArgs(
				teamID,
				countryCode,
				"firstname",
				"lastname",
				int32(rand.Intn(maxAge-minAge+1)+minAge),
				position,
				decimal.NewFromInt(price),
			)
			players = append(players, player)
		}
	}
	return players
}
