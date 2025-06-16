package user

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	h := newHandler(
		c.Services.UserService,
		c.Cfg.Pagination.S,
		c.Cfg.Pagination.M,
	)

	group.GET("", h.List, m.JWTMiddleware, m.IsAdmin)

	meGroup := group.Group("/me")
	
	meGroup.GET("", h.GetMe)
	meGroup.DELETE("", h.DeleteMe)
	meGroup.PUT("/change-password", h.ChangePasswordMe)
}
