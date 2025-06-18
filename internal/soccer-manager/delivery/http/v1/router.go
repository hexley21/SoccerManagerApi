package v1

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/auth"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/globe"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/player"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/player_position"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/team"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/transfer"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery/http/v1/user"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	globe.RegisterRoutes(g.Group("/globe"), c)

	auth.RegisterRoutes(g.Group("/auth"), c)
	user.RegisterRoutes(g.Group("/users"), c, m)

	team.RegisterRoutes(g, c, m)

	player_position.RegisterRoutes(g, c, m)
	player.RegisterRoutes(g, c, m)

	transfer.RegisterRoutes(g, c, m)
}
