package player_position

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	h := NewHandler(
		c.Services.PlayerPosService,
	)

	g.GET("/player-positions", h.ListTranslated, m.AcceptLanguage)
	g.GET("/player-positions/codes", h.ListCodes)

	trG := g.Group("/player-positions/translations")
	trG.GET("", h.ListAllTranslations)
	trG.POST("", h.CreateTranslation, m.JWTMiddleware, m.IsAdmin)
	trG.PUT("", h.UpdateTranslation, m.JWTMiddleware, m.IsAdmin)
	trG.DELETE("", h.DeleteTranslation, m.JWTMiddleware, m.IsAdmin)
}
