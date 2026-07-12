package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"asset-flow/internal/config"
	"asset-flow/internal/handler"
	"asset-flow/internal/repository"
	"asset-flow/internal/router"
	"asset-flow/internal/server"
	"asset-flow/internal/middleware"
	"time"

	"github.com/labstack/echo/v4"
)

const DefaultContextTimeout = 30

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(
			"failed to load config: ",
			err,
		)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal(
			"failed to initialize server: ",
			err,
		)
	}

	repo := repository.New(
		srv.DB.Pool,
	)

	h := handler.New(
		repo,
	)

	e := echo.New()

	e.Use(
		middleware.CORS(cfg),
		middleware.Logger(),
	)

	router.Register(
		e,
		h,
	)

	srv.SetupHTTPServer(e)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)

	go func() {

		if err := srv.Start(); err != nil &&
			!errors.Is(
				err,
				http.ErrServerClosed,
			) {

			log.Fatal(
				"failed to start server: ",
				err,
			)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel :=
		context.WithTimeout(
			context.Background(),
			DefaultContextTimeout*time.Second,
		)

	defer cancel()
	defer stop()

	if err := srv.Shutdown(
		shutdownCtx,
	); err != nil {

		log.Fatal(
			"server forced to shutdown: ",
			err,
		)
	}

	log.Println(
		"server exited properly",
	)
}