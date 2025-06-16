package team

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/delivery"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(group *echo.Group, c *delivery.Components, m *delivery.Middlewares) {
	h := newHandler(c.Services.TeamService, c.Cfg.Pagination.S, c.Cfg.Pagination.M)

	group.GET("/teams", h.GetTeams, m.AcceptLanguage)
	group.GET("/teams/:team_id", h.GetTeamById, m.AcceptLanguage)
	group.GET("/teams/:team_id/locales", h.GetAvailableLocales)
	group.GET("/users/:user_id/team", h.GetTeamByUserId, m.AcceptLanguage)

	
	selfGroup := group.Group("/teams/me")
	
	selfGroup.PUT("/update-country", h.UpdateTeamCountry, m.JWTMiddleware)
	selfGroup.GET("/translations", h.GetSelfTeamTranslations, m.JWTMiddleware)
	selfGroup.POST("/translations", h.CreateSelfTeamTranslation, m.JWTMiddleware)
	selfGroup.PUT("/translations", h.UpdateSelfTeamTranslation, m.JWTMiddleware)
	selfGroup.DELETE("/translations", h.DeleteSelfTeamTranslation, m.JWTMiddleware)
}
