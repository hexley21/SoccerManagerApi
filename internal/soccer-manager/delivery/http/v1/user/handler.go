package user

import (
	"errors"
	"net/http"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	userService service.UserService
	pageSize    int32
	pageLimit   int32
}

func newHandler(userService service.UserService, pageSize int32, pageLimit int32) *handler {
	return &handler{
		userService: userService,
		pageSize:    pageSize,
		pageLimit:   pageLimit,
	}
}

// @Summary List users (ADMIN)
// @Description Get a paginated list of users
// @Tags users
// @Security AccessToken
// @Produce json
// @Param cursor query int false "Pagination cursor"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} common.apiResponse{data=[]userResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/users [get]
func (h *handler) List(c echo.Context) error {
	// parse pagination params
	pagination, err := common.ParsePagination(c, h.pageSize, h.pageLimit)
	if err != nil {
		return err
	}

	users, err := h.userService.List(c.Request().Context(), pagination.Cursor, pagination.PageSize)
	if err != nil {
		c.Logger().Errorf("failed to fetch user list: %v", err)
		return err
	}

	res := make([]userResponseDTO, len(users))
	for i, usr := range users {
		res[i] = NewUserResponseDTO(usr.ID, usr.Username, string(usr.Role))
	}

	return c.JSON(http.StatusOK, common.NewApiResponse(res))
}

// @Summary Get current user
// @Description Get the authenticated user's data
// @Tags users
// @Security AccessToken
// @Produce json
// @Success 200 {object} common.apiResponse{data=userResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/users/me [get]
func (h *handler) GetMe(c echo.Context) error {
	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(access.NewInvalidTokenError(userData))
	}

	user, err := h.userService.Get(c.Request().Context(), userData.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.JSON(
		http.StatusOK,
		common.NewApiResponse(NewUserResponseDTO(user.ID, user.Username, string(user.Role))),
	)
}

// @Summary Delete current user
// @Description Delete the authenticated user's account
// @Tags users
// @Security AccessToken
// @Success 204 "No Content"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/users/me [delete]
func (h *handler) DeleteMe(c echo.Context) error {
	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		return echo.ErrUnauthorized.WithInternal(access.NewInvalidTokenError(userData))
	}

	if err := h.userService.Delete(c.Request().Context(), userData.UserID); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return echo.ErrNotFound.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary Change current user's password
// @Description Change the authenticated user's password
// @Tags users
// @Accept json
// @Security AccessToken
// @Param request body updatePasswordRequestDTO true "Password update details"
// @Success 200 "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/users/me/change-password [put]
func (h *handler) ChangePasswordMe(c echo.Context) error {
	var req updatePasswordRequestDTO
	if err := c.Bind(&req); err != nil {
		c.Logger().Debug(err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Debug(err)
		return echo.ErrBadRequest.WithInternal(err)
	}

	userData, ok := c.Get(access.CtxKey).(access.Data)
	if !ok {
		err := access.NewInvalidTokenError(userData)
		c.Logger().Error(err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	if err := h.userService.UpdatePassword(c.Request().Context(), userData.UserID, req.OldPassword, req.NewPassword); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.Logger().Error("failed to find user on password-change: %v", err)
			return echo.ErrNotFound.WithInternal(err)
		}
		if errors.Is(err, service.ErrIncorrectPassword) {
			return echo.ErrUnauthorized.WithInternal(err)
		}

		return err
	}

	return c.NoContent(http.StatusOK)
}
