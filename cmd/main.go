package main

import (
	"db_monitoring_daemon/internal/client"
	"db_monitoring_daemon/internal/daemon"
	"db_monitoring_daemon/internal/database"
	"db_monitoring_daemon/internal/server"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// "db_monitoring_daemon/int"

func JsonParser(path string, config interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(byteValue, config)
}

func main() {
	now := time.Now()
	f, err := os.Create(fmt.Sprintf("log/log_%s.log", now.Format(time.RFC3339)))
	if err != nil {
		panic("cannot open log file")
	}
	defer f.Close()
	logger := zerolog.New(f).With().Timestamp().Logger()
	logger.Level(zerolog.DebugLevel)
	log.Output(f)

	databaseConfig := &database.DatabaseConfig{}
	err = JsonParser("configs/database.json", databaseConfig)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	database := database.NewDatabase(databaseConfig, logger, "internal/sql/")
	database.SetRequest("drop_table")
	database.SetRequest("checkpoint")
	database.SetRequest("select_count_star")
	database.SetRequest("truncate_table")
	database.SetRequest("pg_terminate_backend")
	database.SetRequest("select_pg_stat_activity")

	db := database.Connect()
	if db == nil {
		logger.Fatal().Msg("cannot connect to the database")
		return
	}
	db.Close()

	monitoringCfg := &client.ServiceConfig{}
	err = JsonParser("configs/monitoring.json", monitoringCfg)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	monitoring := client.NewMonitoring(monitoringCfg, logger)

	metricaCfg := &client.ServiceConfig{}
	err = JsonParser("configs/monitoring.json", monitoringCfg)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	metrica := client.NewMetrica(metricaCfg, logger)

	daemonConfig := &daemon.DaemonConfig{}
	err = JsonParser("configs/daemon.json", daemonConfig)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	daemon := daemon.NewDaemon(
		daemonConfig,
		monitoring,
		metrica,
		logger,
		database,
	)

	serverConfig := &server.ServerConfig{}
	err = JsonParser("configs/server.json", serverConfig)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	server := server.NewServer(serverConfig, logger, database)

	go daemon.Start()
	server.Start()
}
