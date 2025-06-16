package team

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

// TODO: add team assignation on register
type handler struct {
	teamService service.TeamService
	pageSize    int32
	pageLimit   int32
}

func newHandler(teamService service.TeamService, pageSize int32, pageLimit int32) *handler {
	return &handler{
		teamService: teamService,
		pageSize:    pageSize,
		pageLimit:   pageLimit,
	}
}

// @Summary List teams
// @Description Returns a paginated list of teams
// @Tags team
// @Produce json
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]teamResponseDTO} "OK"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams [get]
func (h *handler) GetTeams(c echo.Context) error {
	locale, ok := c.Get(domain.LocaleCtxKey).(domain.LocaleCode)
	if !ok {
		c.Logger().Debug("couldn't get locale when fetching teams")
		locale = domain.LocaleCode("")
	}

	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	teams, err := h.teamService.GetTeams(
		c.Request().Context(),
		locale,
		pagination.Cursor,
		pagination.PageSize,
	)
	if err != nil {
		return err
	}

	res := make([]teamResponseDTO, len(teams))
	for i, t := range teams {
		res[i] = TeamResponseAdapter(t)
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Get team by ID
// @Description Returns a single team's data by ID
// @Tags team
// @Produce json
// @Param team_id path int true "Team ID"
// @Success 200 {object} common.apiResponse{data=teamResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/{team_id} [get]
func (h *handler) GetTeamById(c echo.Context) error {
	teanID, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	locale, ok := c.Get(domain.LocaleCtxKey).(domain.LocaleCode)
	if !ok {
		c.Logger().Debug("couldn't get locale when fetching teams")
		locale = domain.LocaleCode("")
	}

	team, err := h.teamService.GetTeamById(
		c.Request().Context(),
		locale,
		teanID,
	)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(TeamResponseAdapter(team)))
}

// @Summary Get team by user ID
// @Description Returns a single team by user ID
// @Tags team
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} common.apiResponse{data=teamResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/users/{user_id}/team [get]
func (h *handler) GetTeamByUserId(c echo.Context) error {
	userId, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	locale, ok := c.Get(domain.LocaleCtxKey).(domain.LocaleCode)
	if !ok {
		c.Logger().Debug("couldn't get locale when fetching teams")
		locale = domain.LocaleCode("")
	}

	team, err := h.teamService.GetTeamByUserId(
		c.Request().Context(),
		locale,
		userId,
	)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(TeamResponseAdapter(team)))
}

// @Summary Update team's country
// @Description Updates the country of the authenticated user's team
// @Tags team
// @Accept json
// @Produce json
// @Security AccessToken
// @Param request body updateTeamCountryRequestDTO true "Country code"
// @Success 200 {object} echo.HTTPError "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/me/update-country [put]
func (h *handler) UpdateTeamCountry(c echo.Context) error {
	var req updateTeamCountryRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		err := access.NewInvalidTokenError(userData)
		c.Logger().Error(err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	err := h.teamService.UpdateTeamCountry(
		c.Request().Context(),
		userData.UserID,
		domain.CountryCode(req.CountryCode),
	)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		if errors.Is(err, service.ErrNonexistentCode) {
			return echo.ErrBadRequest.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}

// @Summary List translations
// @Description Returns translations for the authenticated user's team
// @Tags team
// @Produce json
// @Security AccessToken
// @Success 200 {object} common.apiResponse{data=[]domain.TeamTranslation} "OK"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/me/translations [get]
func (h *handler) GetSelfTeamTranslations(c echo.Context) error {
	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		err := access.NewInvalidTokenError(userData)
		c.Logger().Error(err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	translations, err := h.teamService.GetTeamTranslations(c.Request().Context(), userData.UserID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(translations))
}

// @Summary Create team translation
// @Description Creates a new translation for the authenticated user's team
// @Tags team
// @Accept json
// @Produce json
// @Security AccessToken
// @Param request body createTeamTranslationRequestDTO true "Translation details"
// @Success 201 "Created"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 409 {object} echo.HTTPError "Conflict"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/me/translations [post]
func (h *handler) CreateSelfTeamTranslation(c echo.Context) error {
	var req createTeamTranslationRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		err := access.NewInvalidTokenError(userData)
		c.Logger().Error(err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	if err := h.teamService.CreateTeamTranslation(c.Request().Context(), req.Locale, userData.UserID, req.Name); err != nil {
		if errors.Is(err, service.ErrTranslationExists) {
			return echo.ErrConflict.WithInternal(err)
		}
		if errors.Is(err, service.ErrNonexistentCode) {
			return echo.ErrBadRequest.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusCreated)
}

// @Summary Update team translation
// @Description Updates an existing translation for the authenticated user's team
// @Tags team
// @Accept json
// @Produce json
// @Security AccessToken
// @Param request body updateTeamTranslationRequestDTO true "Translation details"
// @Success 204 "No Content"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/me/translations [put]
func (h *handler) UpdateSelfTeamTranslation(c echo.Context) error {
	var req updateTeamTranslationRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		err := access.NewInvalidTokenError(userData)
		c.Logger().Error(err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	if err := h.teamService.UpdateTeamTranslation(c.Request().Context(), req.Locale, userData.UserID, req.Name); err != nil {
		if errors.Is(err, service.ErrTranslationNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary Delete team translation
// @Description Deletes a team translation for the authenticated user's team
// @Tags team
// @Produce json
// @Security AccessToken
// @Param request body deleteTeamTranslationRequestDTO true "Delete details"
// @Success 204 "No Content"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/me/translations [delete]
func (h *handler) DeleteSelfTeamTranslation(c echo.Context) error {
	var req deleteTeamTranslationRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		err := access.NewInvalidTokenError(userData)
		c.Logger().Error(err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	if err := h.teamService.DeleteTeamTranslation(c.Request().Context(), req.Locale, userData.UserID); err != nil {
		if errors.Is(err, service.ErrTranslationNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary Get available locales for team by ID
// @Description Returns a list of available translation locales
// @Tags team
// @Produce json
// @Param team_id path int true "Team ID"
// @Success 200 {object} common.apiResponse{data=teamResponseDTO} "OK"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/teams/{team_id}/locales [get]
func (h *handler) GetAvailableLocales(c echo.Context) error {
	teamId, err := strconv.ParseInt(c.Param("team_id"), 10, 64)
	if err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	locales, err := h.teamService.GetAvailableLocales(
		c.Request().Context(),
		teamId,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(locales))
}
