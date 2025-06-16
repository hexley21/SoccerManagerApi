package middleware

import (
	"net/http"
	"strings"

	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/labstack/echo/v4"
)

// JWTAuth is a middleware that checks for JWT token in the request and validates it.
// If the token is not present or invalid, it returns 401 Unauthorized.
// If the token is valid, it adds user data to the request context.
//
// It skips the middleware for all GET and non administrative requests.
func JWTAuth(jwtManager access.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// get token from header
			token := c.Request().Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			if token == "" {
				return JSONErr(c, http.StatusUnauthorized, ErrAuthHeaderRequired)
			}

			// parse token
			externalUser, err := jwtManager.ParseTokenString(token)
			if err != nil {
				return JSONErr(c, http.StatusUnauthorized, ErrInvalidToken)
			}

			// add user data to context
			c.Set(access.CtxKey, externalUser)

			return next(c)
		}
	}
}
