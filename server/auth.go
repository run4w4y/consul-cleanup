package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthStaticToken(token string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			authFields := strings.Fields(authHeader)

			if len(authFields) != 2 || !strings.EqualFold(authFields[0], "bearer") || authFields[1] != token {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			return next(c)
		}
	}
}
