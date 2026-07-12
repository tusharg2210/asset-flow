package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"asset-flow/internal/db"
	"asset-flow/internal/config"
)

type Server struct {
	Config *config.Config

	DB *db.DB

	httpServer *http.Server
}

func New(
	cfg *config.Config,
) (*Server, error) {

	database, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to initialize database: %w",
			err,
		)
	}

	server := &Server{
		Config: cfg,
		DB:     database,
	}

	return server, nil
}

func (s *Server) SetupHTTPServer(
	handler http.Handler,
) {

	s.httpServer = &http.Server{
		Addr: ":" + s.Config.Server.Port,

		Handler: handler,

		ReadTimeout: time.Duration(
			s.Config.Server.ReadTimeout,
		) * time.Second,

		WriteTimeout: time.Duration(
			s.Config.Server.WriteTimeout,
		) * time.Second,

		IdleTimeout: time.Duration(
			s.Config.Server.IdleTimeout,
		) * time.Second,
	}
}

func (s *Server) Start() error {

	if s.httpServer == nil {
		return errors.New(
			"http server not initialized",
		)
	}

	log.Printf(
		"starting server on port %s (%s)",
		s.Config.Server.Port,
		s.Config.Primary.Env,
	)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(
	ctx context.Context,
) error {

	if err := s.httpServer.Shutdown(
		ctx,
	); err != nil {

		return fmt.Errorf(
			"failed to shutdown http server: %w",
			err,
		)
	}

	s.DB.Pool.Close()

	return nil
}
