package v1

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/globe"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	globe.RegisterRoutes(group.Group("/globe"), c)
}
