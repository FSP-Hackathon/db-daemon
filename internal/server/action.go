package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type RequestBody struct {
	Pid       int64  `json:"pid"`
	TableName string `json:"table_name"`
}

func (s *Server) checkpoint() {
	s.logger.Debug().Msg("checkpoint")
	rows := s.database.ExeсRequest("checkpoint")
	if rows == nil {
		log.Warn().Msg("cannot execute command checkpoint")
		return
	}
	defer rows.Close()
}

func (s *Server) terminate(pid int64) {
	s.logger.Debug().Msg("terminate")
	rows := s.database.ExeсRequest("pg_terminate_backend", pid)
	if rows == nil {
		log.Warn().Msg("cannot execute command terminate")
		return
	}
	defer rows.Close()
}

func (s *Server) ActionPost(w http.ResponseWriter, r *http.Request) {
	action_type := strings.TrimPrefix(r.URL.Path, "/api/action/")

	var body RequestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if action_type == "checkpoint" {
		s.checkpoint()
	} else if action_type == "terminate" {
		s.terminate(body.Pid)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
