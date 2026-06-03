package middleware

import (
	myjwt "my-project/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(handler echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"missing token",
			)
		}

		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		_, err := myjwt.ValidateToken(tokenString)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusUnauthorized,
				"invalid token",
			)
		}

		return handler(c)
	}
}
