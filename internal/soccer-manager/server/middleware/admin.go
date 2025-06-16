package middleware

import (
	"net/http"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/labstack/echo/v4"
)

func IsAdmin() echo.MiddlewareFunc {
	return isRole(domain.UserRoleADMIN)
}

func IsUser() echo.MiddlewareFunc {
	return isRole(domain.UserRoleUSER)
}

// isRole just checks if the role from JWTAuth's access key matches with the provided role
// Use for placing barriers between admin and user
func isRole(role domain.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userData, ok := c.Get(access.CtxKey).(access.Data)
			if !ok {
				c.Logger().Error("invalid access token: %v", userData)
				return JSONErr(c, http.StatusForbidden, ErrInvalidToken)
			}

			if userData.Role != role {
				c.Logger().Error("incorrect role: had - %v, wanted - %v", userData.Role, role)
				return JSONErr(c, http.StatusForbidden, ErrInsufficientRights)
			}

			return next(c)
		}
	}
}
