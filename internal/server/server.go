package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type Server struct {
	cfg    *ServerConfig
	logger zerolog.Logger
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewServer(cfg *ServerConfig, logger zerolog.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Start() {
	s.logger.Info().Msg("start server")
	http.HandleFunc("/action", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("lol, action")
	})
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port), nil)
}
