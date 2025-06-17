package player

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components) {
	h := NewHandler(
		c.Services.PlayerService,
		c.Cfg.Pagination.S,
		c.Cfg.Pagination.M,
	)

	g.GET("/players", h.GetAllPlayers)
	g.GET("/players/:player_id", h.GetPlayerById)
	g.PUT("/players", h.UpdatePlayerData)
	g.GET("/teams/:team_id/players", h.GetPlayersByTeamId)
	g.GET("/users/:user_id/players", h.GetPlayersByUserId)
}
