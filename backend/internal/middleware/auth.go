package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/util"
)

func RequireAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return echo.NewHTTPError(http.StatusUnauthorized, "authorization header must be 'Bearer <token>'")
			}

			claims, err := util.ParseAccessToken(jwtSecret, parts[1])
			if err != nil {
				if errors.Is(err, util.ErrWrongTokenType) {
					return echo.NewHTTPError(http.StatusUnauthorized, "refresh token cannot be used here")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)
			return next(c)
		}
	}
}

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("user_role").(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthenticated")
			}
			if _, ok := allowed[role]; !ok {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}
			return next(c)
		}
	}
}