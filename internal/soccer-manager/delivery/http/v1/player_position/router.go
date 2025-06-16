package player_position

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	h := NewHandler(
		c.Services.PlayerPosService,
	)

	group.GET("", h.ListTranslated, m.AcceptLanguage)
	group.GET("/codes", h.ListCodes)

	group.GET("/translations", h.ListAllTranslations)
	group.POST("/translations", h.CreateTranslation, m.JWTMiddleware, m.IsAdmin)
	group.PUT("/translations", h.UpdateTranslation, m.JWTMiddleware, m.IsAdmin)
	group.DELETE("/translations", h.DeleteTranslation, m.JWTMiddleware, m.IsAdmin)
}
