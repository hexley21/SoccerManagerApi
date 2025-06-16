package middleware

import "github.com/labstack/echo/v4"

const (
	ErrAuthHeaderRequired = "authorization header is required"
	ErrInvalidToken       = "invalid token"
	ErrInsufficientRights = "insufficient rights"
)

func JSONErr(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"message": message})
}
