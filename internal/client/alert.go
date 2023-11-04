package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AlertBody struct {
	Msg string `json:"msg"`
}

func (c *Monitoring) SendAlert(body AlertBody) {
	c.logger.Debug().Msg("Send allert with message: " + body.Msg)
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		c.logger.Warn().Msg("cannot marshal alert body")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/alert", c.cfg.Host, c.cfg.Port), bytes.NewReader(bodyBytes))
	if err != nil {
		c.logger.Warn().Msg("cannot send alert request: " + err.Error())
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", c.cfg.Token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.logger.Warn().Msg(fmt.Sprintf("error sending alert: %s", err.Error()))
		return
	}
	if res.StatusCode != 200 {
		c.logger.Warn().Msg(fmt.Sprintf("sending alert not ok: %d", res.StatusCode))
		return
	}
}
