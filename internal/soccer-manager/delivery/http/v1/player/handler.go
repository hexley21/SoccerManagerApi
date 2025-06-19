package player

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	playerService service.PlayerService
	pageSize      int32
	pageLimit     int32
}

func newHandler(playerService service.PlayerService, pageSize int32, pageLimit int32) *handler {
	return &handler{
		playerService: playerService,
		pageSize:      pageSize,
		pageLimit:     pageLimit,
	}
}

// @Summary Get all players
// @Description Returns a paginated list of all players
// @Tags players
// @Produce json
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]playerResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/players [get]
func (h *handler) GetAllPlayers(c echo.Context) error {
	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	players, err := h.playerService.GetAllPlayers(
		c.Request().Context(),
		pagination.Cursor,
		pagination.PageSize,
	)
	if err != nil {
		return err
	}

	res := make([]playerResponseDTO, len(players))
	for i, pl := range players {
		res[i] = playerResponseAdapter(pl)
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Get player by ID
// @Description Returns a single player's data by ID
// @Tags players
// @Produce json
// @Param player_id path int true "Player ID"
// @Success 200 {object} common.apiResponse{data=playerResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/players/{player_id} [get]
func (h *handler) GetPlayerById(c echo.Context) error {
	playerId, err := strconv.ParseInt(c.Param("player_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	player, err := h.playerService.GetPlayerById(
		c.Request().Context(),
		playerId,
	)
	if err != nil {
		if errors.Is(err, service.ErrPlayerNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		
		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(playerResponseAdapter(player)))
}

// @Summary Update player data
// @Description Updates a player's data by ID
// @Tags players
// @Accept json
// @Produce json
// @Param player_id path int true "Player ID"
// @Security AccessToken
// @Param request body updatePlayerDataRequestDTO true "Player data to update"
// @Success 200 "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/players/{player_id} [put]
func (h *handler) UpdatePlayerData(c echo.Context) error {
	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(access.NewInvalidTokenError(userData))
	}

	c.Logger().Debug(userData)

	playerId, err := strconv.ParseInt(c.Param("player_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	var req updatePlayerDataRequestDTO
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	err = h.playerService.UpdatePlayerData(
		c.Request().Context(),
		userData.UserID,
		playerId,
		req.FirstName,
		req.LastName,
		req.CountryCode,
	)
	if err != nil {
		if errors.Is(err, service.ErrPlayerNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Get players by team ID
// @Description Returns a paginated list of players for the specified team
// @Tags players
// @Produce json
// @Param team_id path int true "Team ID"
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]playerResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/{team_id}/players [get]
func (h *handler) GetPlayersByTeamId(c echo.Context) error {
	teamId, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	players, err := h.playerService.GetPlayersByTeamId(
		c.Request().Context(),
		teamId,
		pagination.Cursor,
		pagination.PageSize,
	)
	if err != nil {
		return err
	}

	res := make([]playerResponseDTO, len(players))
	for i, pl := range players {
		res[i] = playerResponseAdapter(pl)
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Get players by user ID
// @Description Returns a paginated list of players for the specified user
// @Tags players
// @Produce json
// @Param user_id path int true "User ID"
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]playerResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/users/{user_id}/players [get]
func (h *handler) GetPlayersByUserId(c echo.Context) error {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	players, err := h.playerService.GetPlayersByUserId(
		c.Request().Context(),
		userId,
		pagination.Cursor,
		pagination.PageSize,
	)
	if err != nil {
		return err
	}

	res := make([]playerResponseDTO, len(players))
	for i, pl := range players {
		res[i] = playerResponseAdapter(pl)
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}
