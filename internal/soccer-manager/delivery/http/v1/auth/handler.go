package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/hexley21/soccer-manager/internal/common"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/refresh"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/service"
	"github.com/labstack/echo/v4"
)

const refreshCookie = "refresh_token"

type handler struct {
	authService       service.AuthService
	accessJWTManager  access.Manager
	refreshJWTManager refresh.Manager
}

func NewHandler(
	authService service.AuthService,
	accessJWTManager access.Manager,
	refreshJWTManager refresh.Manager,
) *handler {
	return &handler{
		authService:       authService,
		accessJWTManager:  accessJWTManager,
		refreshJWTManager: refreshJWTManager,
	}
}

// @Summary Login user
// @Description Authenticates user with username and password and returns JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequestDTO true "Login credentials"
// @Success 200 {object} common.apiResponse{data=loginResponseDTO} "OK"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 404 {object} echo.HTTPError "Not Found"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/auth/login [post]
func (h *handler) Login(c echo.Context) error {
	var req loginRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, err := h.authService.Authenticate(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.Logger().Debug(err)
			return echo.ErrNotFound.WithInternal(err)
		}
		if errors.Is(err, service.ErrIncorrectPassword) {
			c.Logger().Debug(err)
			return echo.ErrUnauthorized.WithInternal(err)
		}

		c.Logger().Error("failed to login user: %v", err)
		return err
	}

	accessToken, err := h.accessJWTManager.CreateTokenString(
		access.NewData(user.ID, user.Role),
	)
	if err != nil {
		c.Logger().Errorf("failed to create access_token: %v", err)
		return err
	}

	refreshToken, err := h.refreshJWTManager.CreateTokenString(
		refresh.NewData(user.ID),
	)
	if err != nil {
		c.Logger().Errorf("failed to create refresh_token: %v", err)
		return err
	}

	setCookieJWT(c, refreshCookie, refreshToken, h.refreshJWTManager.TTL())

	return c.JSON(http.StatusOK, common.NewApiResponse(loginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // the refresh token is left only for better apie usage
	}))
}

// @Summary Logout user
// @Description Logs out user by invalidating refresh token cookie
// @Tags auth
// @Success 204 "No Content"
// @Router /v1/auth/logout [post]
func (h *handler) Logout(c echo.Context) error {
	eraseCookie(c, refreshCookie)
	return c.NoContent(http.StatusNoContent)
}

// @Summary Register new user
// @Description Creates a new user account and returns JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequestDTO true "Registration details"
// @Success 201 {object} common.apiResponse{data=registerResponseDTO} "Created"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Failure 409 {object} echo.HTTPError "Conflict"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/auth/register [post]
func (h *handler) Register(c echo.Context) error {
	var req registerRequestDTO
	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	user, err := h.authService.CreateUser(
		c.Request().Context(),
		req.Username,
		req.Password,
		req.Role,
	)
	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) {
			c.Logger().Debug(err)
			return echo.ErrConflict.WithInternal(err)
		}

		c.Logger().Errorf("failed to create user: %v", err)
		return err
	}

	accessToken, err := h.accessJWTManager.CreateTokenString(
		access.NewData(user.ID, user.Role),
	)
	if err != nil {
		c.Logger().Errorf("failed to create access_token: %v", err)
		return err
	}

	refreshToken, err := h.refreshJWTManager.CreateTokenString(
		refresh.NewData(user.ID),
	)
	if err != nil {
		c.Logger().Errorf("failed to create refresh_token: %v", err)
		return err
	}

	setCookieJWT(c, refreshCookie, refreshToken, h.refreshJWTManager.TTL())

	return c.JSON(http.StatusCreated, common.NewApiResponse(registerResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}))
}

// @Summary Refresh access token
// @Description Get a new access token using refresh token from cookie
// @Tags auth
// @Produce json
// @Success 200 {object} common.apiResponse{data=refreshResponseDTO} "OK"
// @Failure 401 {object} echo.HTTPError "Unauthorized"
// @Failure 500 {object} echo.HTTPError "Internal Server Error"
// @Router /v1/auth/refresh [post]
func (h *handler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie(refreshCookie)
	if err != nil {
		c.Logger().Errorf("failed to get refresh_token cookie: %v", err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	refreshData, err := h.refreshJWTManager.ParseTokenString(cookie.Value)
	if err != nil {
		c.Logger().Errorf("failed to invalidate refresh_token: %v", err)
		return echo.ErrUnauthorized.WithInternal(err)
	}

	user, err := h.authService.GetUserById(c.Request().Context(), refreshData.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			eraseCookie(c, refreshCookie)
			c.Logger().Errorf("failed to find user from a refresh_token: %v", err)
			return echo.ErrUnauthorized.WithInternal(err)
		}

		return err
	}

	accessToken, err := h.accessJWTManager.CreateTokenString(
		access.NewData(refreshData.UserID, user.Role),
	)
	if err != nil {
		c.Logger().Errorf("failed to sign access_token")
		return err
	}

	return c.JSON(
		http.StatusOK,
		common.NewApiResponse(RefreshResponseDTO{AccessToken: accessToken}),
	)
}

func setCookieJWT(c echo.Context, cookieName string, token string, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     cookieName,
		Value:    token,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(ttl),
	})
}

func eraseCookie(c echo.Context, cookieName string) {
	c.SetCookie(&http.Cookie{
		Name:     cookieName,
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})
}
