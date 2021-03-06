package model

import (
	"encoding/json"

	"github.com/sonod/pisera/config"
	"github.com/sonod/pisera/lib"
)

type Device struct {
	Data []deviceData
}

type deviceData struct {
	Id       string
	Hostname string
	Ip       string
	Rack     string
}

type createDevice struct {
	Sections    string
	Hostname    string
	Ip_addr     string
	Description string
}

func GetDeviceList(conf config.Config, token string) (*Device, error) {
	urlPath := conf.Server.AppID + "/devices/"
	d, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	device := new(Device)
	if err := json.Unmarshal(d, device); err != nil {
		return nil, err
	}

	return device, nil
}

func CreateDevice(conf config.Config, token, secID, hostName, ipmi string) ([]byte, error) {
	urlPath := conf.Server.AppID + "/devices/"
	cdBody := &createDevice{
		Sections:    secID,
		Hostname:    hostName,
		Ip_addr:     ipmi,
		Description: "IPMI",
	}
	res, err := client.PHPIPAMRequest(conf, "POST", token, "", urlPath, cdBody)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateDevice(conf config.Config, token, DevID, ipmi string) ([]byte, error) {
	urlPath := conf.Server.AppID + "/devices/"
	udBody := &deviceData{
		Id: DevID,
		Ip: ipmi,
	}
	res, err := client.PHPIPAMRequest(conf, "PATCH", token, "", urlPath, udBody)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SearchDevices(conf config.Config, token, hostname string) (*Device, error) {
	urlPath := conf.Server.AppID + "/devices/search/" + hostname + "/"
	d, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	device := new(Device)
	if err := json.Unmarshal(d, device); err != nil {
		return nil, err
	}

	return device, nil
}

func DeleteDevice(conf config.Config, token, id string) ([]byte, error) {
	urlPath := conf.Server.AppID + "/devices/"
	ddBody := map[string]interface{}{"id": id}
	res, err := client.PHPIPAMRequest(conf, "DELETE", token, "", urlPath, ddBody)
	if err != nil {
		return nil, err
	}
	return res, nil
}
