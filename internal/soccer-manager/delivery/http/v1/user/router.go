package user

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	h := newHandler(
		c.Services.UserService,
		c.Cfg.Pagination.S,
		c.Cfg.Pagination.M,
	)

	g.GET("", h.List, m.JWTMiddleware, m.IsAdmin)

	meGroup := g.Group("/me", m.JWTMiddleware)

	meGroup.GET("", h.GetMe)
	meGroup.DELETE("", h.DeleteMe)
	meGroup.PUT("/change-password", h.ChangePasswordMe)
}
