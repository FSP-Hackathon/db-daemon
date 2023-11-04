package main

import (
	"db_monitoring_daemon/internal/client"
	"db_monitoring_daemon/internal/daemon"
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
	log.Output(f)

	logger.Warn().Msg("start daemon")

	monitoringCfg := &client.MonitringConfig{}
	err = JsonParser("configs/daemon.json", monitoringCfg)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	monitoring := client.NewMonitoring(monitoringCfg, logger)

	daemonConfig := &daemon.DaemonConfig{}
	err = JsonParser("configs/monitoringclient.json", daemonConfig)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	daemon := daemon.NewDaemon(daemonConfig, monitoring, logger)

	serverConfig := &server.ServerConfig{}
	err = JsonParser("configs/server.json", serverConfig)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	// server := server.NewServer(serverConfig, logger)

	data := daemon.Update()
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(bytes))

	// go daemon.Start()
	// server.Start()
}
