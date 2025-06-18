package transfer

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	h := newHandler(c.Services.TransferService, c.Cfg.Pagination.M, c.Cfg.Pagination.L)

	g.GET("/transfers", h.GetTransfers)
	g.GET("/transfers/:transfer_id", h.GetTransferById)
	g.GET("/teams/:team_id/transfers", h.GetTransfersByTeamId)
	g.GET("/players/:player_id/transfer", h.GetTransferByPlayer)

	g.POST("/transfers", h.CreateTransfer, m.JWTMiddleware)
	g.DELETE("/transfers/:transfer_id", h.DeleteTransfer, m.JWTMiddleware)
	g.PUT("/transfers/:transfer_id", h.UpdateTransfer, m.JWTMiddleware)

	g.POST("/transfers/:transfer_id/buy", h.BuyPlayer, m.JWTMiddleware)
}
