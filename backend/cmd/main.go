package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"

	"asset-flow/internal/config"
	"asset-flow/internal/handler"
	appmw "asset-flow/internal/middleware"
	"asset-flow/internal/repository"
	"asset-flow/internal/router"
	"asset-flow/internal/server"
)

type echoValidator struct {
	v *validator.Validate
}

func (ev *echoValidator) Validate(i interface{}) error {
	return ev.v.Struct(i)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	repos := repository.New(srv.DB.Pool)
	handlers := handler.New(repos, srv.DB.Pool, cfg)

	e := echo.New()
	e.HideBanner = true
	e.Validator = &echoValidator{v: validator.New()}

	e.Use(echomw.Recover())
	e.Use(echomw.RequestID())
	e.Use(appmw.Logger())
	e.Use(appmw.CORS(cfg))

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
	})

	router.Register(e, cfg, handlers)

	srv.SetupHTTPServer(e)

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped cleanly")
}