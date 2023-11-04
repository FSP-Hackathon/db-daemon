package daemon

import (
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
)

type PgStatActivity struct {
	Pid       int64         `json:"pid"`
	WaitEvent string        `json:"wait_event"`
	Duration  time.Duration `json:"duration"`
}

func (d *Daemon) GetPgStatsActivity() []PgStatActivity {
	var pgstats []PgStatActivity
	rows := d.database.ExeсRequest("select_pg_stat_activity")
	if rows == nil {
		log.Warn().Msg("cannot get pg_stat_activity")
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var pid int64
		var wait_event sql.NullString
		var now time.Time
		var start time.Time
		if err := rows.Scan(&pid, &wait_event, &now, &start); err != nil {
			d.logger.Fatal().Msg(err.Error())
			return nil
		}
		if !wait_event.Valid {
			wait_event = sql.NullString{String: ""}
		}

		pgstats = append(pgstats,
			PgStatActivity{
				Pid:       pid,
				WaitEvent: wait_event.String,
				Duration:  now.Sub(start),
			},
		)
	}
	if err := rows.Err(); err != nil {
		d.logger.Fatal().Msg(err.Error())
		return nil
	}

	return pgstats
}
