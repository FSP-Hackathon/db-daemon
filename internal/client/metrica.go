package client

import "github.com/rs/zerolog"

type Metrica struct {
	cfg    *ServiceConfig
	logger zerolog.Logger
}

func NewMetrica(cfg *ServiceConfig, logger zerolog.Logger) *Metrica {
	return &Metrica{
		cfg:    cfg,
		logger: logger,
	}
}
