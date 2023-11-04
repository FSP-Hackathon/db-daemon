package client

import "github.com/rs/zerolog"

type Monitoring struct {
	cfg    *MonitringConfig
	logger zerolog.Logger
}

func NewMonitoring(cfg *MonitringConfig, logger zerolog.Logger) *Monitoring {
	return &Monitoring{
		cfg:    cfg,
		logger: logger,
	}
}

type MonitringConfig struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}
