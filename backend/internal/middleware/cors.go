package middleware

import (
	"asset-flow/internal/config"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func CORS(
	cfg *config.Config,
) echo.MiddlewareFunc {

	return echoMiddleware.CORSWithConfig(
		echoMiddleware.CORSConfig{
			AllowOrigins: cfg.Server.CORSAllowedOrigins,

			AllowMethods: []string{
				"GET",
				"POST",
				"DELETE",
				"OPTIONS",
			},

			AllowHeaders: []string{
				"Content-Type",
				"Authorization",
			},
		},
	)
}