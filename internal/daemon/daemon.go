package daemon

import (
	"db_monitoring_daemon/internal/client"
	"db_monitoring_daemon/internal/database"
	"db_monitoring_daemon/internal/serverstats"
	"encoding/json"
	"fmt"
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
	PgStatActivity []PgStatActivity `json:"pg_stat_activity"`
}

type DaemonConfig struct {
	Period uint64 `json:"period"` // milliseconds
}

type Daemon struct {
	monitoring  *client.Monitoring
	metrica     *client.Metrica
	cfg         *DaemonConfig
	logger      zerolog.Logger
	database    *database.Database
	alertedPids map[int64]bool
}

func NewDaemon(cfg *DaemonConfig, monitoring *client.Monitoring, metrica *client.Metrica, logger zerolog.Logger, database *database.Database) *Daemon {
	return &Daemon{
		cfg:         cfg,
		monitoring:  monitoring,
		metrica:     metrica,
		logger:      logger,
		database:    database,
		alertedPids: make(map[int64]bool),
	}
}

func (d *Daemon) Update() (data DaemonMessage) {
	var (
		cpu  serverstats.CPU
		ram  serverstats.RAM
		disk serverstats.Disk
	)
	data.Hardware.CPU = cpu.GetStats()
	data.Hardware.RAM = ram.GetStats()
	data.Hardware.Disk = disk.GetStats("/")
	data.PgStatActivity = d.GetPgStatsActivity()

	// some analitic
	for i := 0; i < len(data.PgStatActivity); i++ {
		row := data.PgStatActivity[i]
		if row.WaitEvent == "ClientRead" && row.Duration > 15*time.Minute {
			if _, exists := d.alertedPids[row.Pid]; !exists {
				d.monitoring.SendAlert(client.AlertBody{
					Msg: fmt.Sprintf("Session with pid=%d is connected more than 15 minutes", row.Pid),
				})
				d.alertedPids[row.Pid] = true
			}
		}
	}

	// go to the ml

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
