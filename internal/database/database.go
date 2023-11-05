package database

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Database struct {
	cfg    *DatabaseConfig
	logger zerolog.Logger

	requests map[string]string
	path     string
}

func NewDatabase(cfg *DatabaseConfig, logger zerolog.Logger, path string) *Database {
	return &Database{
		cfg:      cfg,
		logger:   logger,
		requests: make(map[string]string),
		path:     path,
	}
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func (d *Database) SetRequest(name string) {
	filename := name + ".sql"
	file, err := os.Open(d.path + filename)
	if err != nil {
		d.logger.Fatal().Msg(fmt.Sprintf("cannot open file: %s", filename))
		return
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		d.logger.Fatal().Msg(fmt.Sprintf("cannot read file: %s", filename))
		return
	}
	d.requests[name] = string(byteValue)
}

func (d *Database) Connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.cfg.Host, d.cfg.Port, d.cfg.User, d.cfg.Password, d.cfg.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		d.logger.Warn().Msg(err.Error())
		return nil
	}

	return db
}

func (d *Database) ExeсRequest(name string, args ...interface{}) *sql.Rows {
	var exists bool
	var request string
	if request, exists = d.requests[name]; !exists {
		log.Warn().Msg(fmt.Sprintf("sql request not found: %s", name))
		return nil
	}

	db := d.Connect()
	if db == nil {
		d.logger.Warn().Msg("cannot open database")
		return nil
	}
	defer db.Close()
	rows, err := db.Query(request, args)
	if err != nil {
		d.logger.Warn().Msg(fmt.Sprintf("cannot execute sql query: %s, error: %s", name, err.Error()))
		return nil
	}

	return rows
}
