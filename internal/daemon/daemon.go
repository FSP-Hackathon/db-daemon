package daemon

import (
	"db_monitoring_daemon/internal/client"
	"db_monitoring_daemon/internal/serverstats"
	"encoding/json"
	"time"

	"github.com/rs/zerolog"
)

// get stats and send it with socket

type DaemonMessage struct {
	Hardware struct {
		CPU  serverstats.CPUStats  `json:"cpu"`
		Disk serverstats.DiskStats `json:"disk"`
		RAM  serverstats.RAMStats  `json:"ram"`
	} `json:"hardware"`
	// PgStatActivity struct {
	// } `json:"PgStatActivity"`
}

type DaemonConfig struct {
	Period uint64 `json:"period"` // milliseconds
}

type Daemon struct {
	client *client.Monitoring
	cfg    *DaemonConfig
	logger zerolog.Logger
}

func NewDaemon(cfg *DaemonConfig, client *client.Monitoring, logger zerolog.Logger) *Daemon {
	return &Daemon{
		cfg:    cfg,
		client: client,
		logger: logger,
	}
}

func (Daemon) Update() (data DaemonMessage) {
	var (
		cpu  serverstats.CPU
		ram  serverstats.RAM
		disk serverstats.Disk
	)
	data.Hardware.CPU = cpu.GetStats()
	data.Hardware.RAM = ram.GetStats()
	data.Hardware.Disk = disk.GetStats("/")

	// check
	// some analitic

	return data
}

func (d *Daemon) Start() {
	d.logger.Warn().Msg("start daemon")
	go func() {
		for {
			start := time.Now()
			data := d.Update()

			bytes, err := json.Marshal(data)
			if err != nil {
				d.logger.Warn().Msg("cannot marshal daemon message")
			}
			d.logger.Info().Msg(string(bytes))
			// send

			end := time.Now()
			// sleep for compliance with the period
			time.Sleep(end.Sub(start) - time.Duration(d.cfg.Period)*time.Millisecond)
		}
	}()
}
