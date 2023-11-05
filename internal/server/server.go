package server

import (
	"db_monitoring_daemon/internal/database"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type Server struct {
	cfg      *ServerConfig
	logger   zerolog.Logger
	database *database.Database
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewServer(cfg *ServerConfig, logger zerolog.Logger, database *database.Database) *Server {
	return &Server{
		cfg:      cfg,
		logger:   logger,
		database: database,
	}
}

func (s *Server) Start() {
	s.logger.Info().Msg("start server")
	http.HandleFunc("/api/action/", s.ActionPost)
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port), nil)
}
