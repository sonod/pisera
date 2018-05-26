package model

import (
	"encoding/json"

	"github.com/sonod/pisera/config"
	"github.com/sonod/pisera/lib"
)

type Token struct {
	Code    int
	Success bool
	Time    float64
	Data    struct {
		Token   string
		Expires string
	}
}

func GetAuthenticationToken(conf config.Config) (string, error) {
	urlPath := conf.Server.AppID + "/user/"
	t, err := client.PHPIPAMRequest(conf, "POST", "", conf.Server.Username, urlPath, nil)
	if err != nil {
		return "", err
	}

	token := new(Token)
	if err := json.Unmarshal(t, token); err != nil {
		return "", err
	}

	return token.Data.Token, nil
}
