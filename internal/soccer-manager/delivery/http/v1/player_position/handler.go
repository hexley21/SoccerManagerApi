package player_position

import (
	"errors"
	"net/http"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	playerPosService service.PlayerPositionService
}

func NewHandler(
	playerPosService service.PlayerPositionService,
) *handler {
	return &handler{
		playerPosService: playerPosService,
	}
}

// @Summary List position codes
// @Description Returns a list of all available player position codes
// @Tags player-positions
// @Produce json
// @Success 200 {object} common.apiResponse{data=[]string} "OK"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/player-positions/codes [get]
func (h *handler) ListCodes(c echo.Context) error {
	codes, err := h.playerPosService.ListPositionCodes(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(codes))
}

// @Summary List player positions (localized)
// @Description Returns a list of player positions, optionally translated by locale
// @Tags player-positions
// @Produce json
// @Success 200 {object} common.apiResponse{data=playerPositionsResponseDTO} "OK"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/player-positions [get]
func (h *handler) ListTranslated(c echo.Context) error {
	var positions []domain.PlayerPosition
	var err error

	locale, ok := c.Get(domain.LocaleCtxKey).(domain.LocaleCode)
	if !ok || locale == "en" {
		positions, err = h.playerPosService.ListPositions(c.Request().Context())
	} else {
		positions, err = h.playerPosService.ListPositionTranslationsByLocale(
			c.Request().Context(),
			locale,
		)
	}

	if err != nil {
		return err
	}

	res := make(playerPositionsResponseDTO, len(positions))
	for i := range positions {
		pos := positions[i]
		res[pos.Code] = pos.Label
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary List all translations
// @Description Returns all translations of player positions for all locales
// @Tags player-positions
// @Produce json
// @Success 200 {object} common.apiResponse{data=playerPositionsWithLocalesResponseDTO} "OK"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/player-positions/translations [get]
func (h *handler) ListAllTranslations(c echo.Context) error {
	positions, err := h.playerPosService.ListPositionTranslations(c.Request().Context())
	if err != nil {
		return err
	}

	res := make(playerPositionsWithLocalesResponseDTO, len(positions))
	for _, pos := range positions {
		if _, ok := res[pos.Locale]; !ok {
			res[pos.Locale] = make(map[domain.PlayerPositionCode]string)
		}
		res[pos.Locale][pos.Code] = pos.Label
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Create translation
// @Description Creates a new player position translation for a given code and locale
// @Tags player-positions
// @Accept json
// @Produce json
// @Security AccessToken
// @Param request body createPlayerPositionTranslationRequestDTO true "Translation details"
// @Success 201 "Created"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 409 {object} echo.HTTPError "Conflict"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/player-positions/translations [post]
func (h *handler) CreateTranslation(c echo.Context) error {
	var req createPlayerPositionTranslationRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.playerPosService.CreateTranslation(c.Request().Context(), req.Code, req.Locale, req.Label); err != nil {
		if errors.Is(err, service.ErrNonexistentCode) {
			return echo.ErrBadRequest.WithInternal(err)
		}
		if errors.Is(err, service.ErrTranslationExists) {
			return echo.ErrConflict.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusCreated)
}

// @Summary Delete translation
// @Description Deletes a specific player position translation by code and locale
// @Tags player-positions
// @Produce json
// @Security AccessToken
// @Param request body deletePlayerPositionTranslationRequestDTO true "Translation delete details"
// @Success 204 "No Content"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/player-positions/translations [delete]
func (h *handler) DeleteTranslation(c echo.Context) error {
	var req deletePlayerPositionTranslationRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.playerPosService.DeleteTranslation(c.Request().Context(), req.Code, req.Locale); err != nil {
		if errors.Is(err, service.ErrTranslationNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary Update translation
// @Description Updates an existing player position translation
// @Tags player-positions
// @Accept json
// @Produce json
// @Security AccessToken
// @Param request body updatePlayerPositionTranslationRequestDTO true "Updated translation details"
// @Success 204 "No Content"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/player-positions/translations [put]
func (h *handler) UpdateTranslation(c echo.Context) error {
	var req updatePlayerPositionTranslationRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := h.playerPosService.UpdateTranslation(c.Request().Context(), req.Code, req.Locale, req.Label); err != nil {
		if errors.Is(err, service.ErrTranslationNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}
