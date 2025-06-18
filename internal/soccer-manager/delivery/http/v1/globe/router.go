package globe

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components) {
	h := NewHandler(
		c.Services.GlobeService,
	)

	g.GET("/locales", h.Locales)
	g.GET("/countries", h.Countries)
}
