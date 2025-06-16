package v1

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/auth"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/globe"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/player_position"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/user"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	globe.RegisterRoutes(group.Group("/globe"), c)

	auth.RegisterRoutes(group.Group("/auth"), c)
	user.RegisterRoutes(group.Group("/users"), c, m)

	player_position.RegisterRoutes(group.Group("/player-positions"), c, m)
}
