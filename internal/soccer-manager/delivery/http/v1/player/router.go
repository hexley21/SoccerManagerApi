package player

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	h := newHandler(
		c.Services.PlayerService,
		c.Cfg.Pagination.M,
		c.Cfg.Pagination.L,
	)

	g.GET("/players", h.GetAllPlayers)
	g.GET("/players/:player_id", h.GetPlayerById)
	g.PUT("/players/:player_id", h.UpdatePlayerData, m.JWTMiddleware)
	g.GET("/teams/:team_id/players", h.GetPlayersByTeamId)
	g.GET("/users/:user_id/players", h.GetPlayersByUserId)
}
