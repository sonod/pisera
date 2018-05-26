package model

import "github.com/sonod/pisera/config"

type createDevice struct {
	Sections    string
	Hostname    string
	Ip_addr     string
	Description string
}

func CreateDevice(conf config.Config, secID, hostName, ipmi string) ([]byte, error) {
	urlPath := conf.Server.AppID + "/devices/"
	cdBody := &createDevice{
		Sections:    secID,
		Hostname:    hostName,
		Ip_addr:     ipmi,
		Description: "IPMI",
	}
	res, err := client.PHPIPAMRequest(conf, "POST", token, "", urlPath, cdbody)
	if err != nil {
		return nil, err
	}
	return res, nil
}
