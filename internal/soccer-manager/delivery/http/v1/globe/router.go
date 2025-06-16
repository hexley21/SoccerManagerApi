package globe

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, c *delivery.Components) {
	h := NewHandler(
		c.Services.GlobeService,
	)

	group.GET("/locales", h.Locales)
	group.GET("/countries", h.Countries)
}
