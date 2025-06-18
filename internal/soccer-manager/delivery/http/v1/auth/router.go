package auth

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components) {
	h := newHandler(
		c.Services.AuthService,
		c.JWTManagers.Access,
		c.JWTManagers.Refresh,
		c.EventBus,
	)

	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout)
	g.POST("/register", h.Register)
	g.POST("/refresh", h.Refresh)
}
