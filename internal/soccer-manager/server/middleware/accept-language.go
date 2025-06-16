package middleware

import (
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
)

// TODO: add quality weight support
func AcceptLanguage() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			acceptLanguage := c.Request().Header.Get("Accept-Language")
			tag, _, err := language.ParseAcceptLanguage(acceptLanguage)
			if err == nil && len(tag) != 0 {
				loc := domain.LocaleCode(tag[0].String())
				if loc.Valid() {
					c.Set(domain.LocaleCtxKey, loc)
				}
			}

			return next(c)
		}
	}
}
