package auth

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, c *delivery.Components) {
	h := newHandler(
		c.Services.AuthService,
		c.JWTManagers.Access,
		c.JWTManagers.Refresh,
		c.EventBus,
	)

	group.POST("/login", h.Login)
	group.POST("/logout", h.Logout)
	group.POST("/register", h.Register)
	group.POST("/refresh", h.Refresh)
}
