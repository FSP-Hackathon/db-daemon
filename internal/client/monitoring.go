package client

import "github.com/rs/zerolog"

type Monitoring struct {
	cfg    *ServiceConfig
	logger zerolog.Logger
}

func NewMonitoring(cfg *ServiceConfig, logger zerolog.Logger) *Monitoring {
	return &Monitoring{
		cfg:    cfg,
		logger: logger,
	}
}

type ServiceConfig struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}
