package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
)

func createEcho(config ServerCleanupConfig, logger zerolog.Logger) *echo.Echo {
	e := echo.New()
	wrappedLogger := lecho.From(logger, lecho.WithLevel(log.Lvl(zerolog.InfoLevel)))
	e.Logger = wrappedLogger

	e.HideBanner = true

	e.Use(middleware.RequestID())
	e.Use(lecho.Middleware(lecho.Config{
		Logger: wrappedLogger,
	}))
	e.Use(middleware.Recover())

	// middleware to use a modified echo context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ApplicationContext{c, config}
			return next(cc)
		}
	})

	api := e.Group("/api/v1")

	api.GET("/health", HealthCheck)

	authorized := api.Group("")
	if config.AccessToken != "" {
		authorized.Use(AuthStaticToken(config.AccessToken))
	}

	authorized.POST("/oneshot", OneshotCleanupAll)
	authorized.POST("/oneshot/:service", OneshotCleanupService)

	return e
}

func CreateEchoWithServer(ctx context.Context, config ServerCleanupConfig) (*echo.Echo, *http.Server) {
	logger := zerolog.Ctx(ctx)

	e := createEcho(config, logger.With().Logger())
	listenAddr := fmt.Sprintf(":%d", config.Port)

	srv := &http.Server{
		Addr:        listenAddr,
		Handler:     e,
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	return e, srv
}
