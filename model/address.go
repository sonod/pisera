package model

import (
	"encoding/json"
	"net"

	"github.com/sonod/pisera/config"
	"github.com/sonod/pisera/lib"
)

type FreeAddress struct {
	Data string
}

type UsedAddresses struct {
	Data []usedAddressesData
}

type usedAddressesData struct {
	IP          string
	Description string
	Hostname    string
	Mac         string
}

type HostAddress struct {
	Data []hostAddressData
}

type hostAddressData struct {
	ID          string
	SubnetID    string
	IP          string
	Description string
	Hostname    string
	Mac         string
}

type createIPAddress struct {
	SubnetId    string `json:"subnetId"`
	IP          string `json:"ip"`
	Hostname    string `json:"hostname"`
	Description string `json:"description"`
	Mac         string `json:"mac"`
}

func GetFirstFreeAddress(conf config.Config, token, snID string) (*FreeAddress, error) {
	urlPath := conf.Server.AppID + "/subnets/" + snID + "/first_free/"
	ffa, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	freeAddress := new(FreeAddress)
	if err := json.Unmarshal(ffa, freeAddress); err != nil {
		return nil, err
	}

	return freeAddress, nil
}

func GetUsedAddressList(conf config.Config, token, snID string) (*UsedAddresses, error) {
	urlPath := conf.Server.AppID + "/subnets/" + snID + "/addresses/"
	ual, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	usedAddresses := new(UsedAddresses)
	if err := json.Unmarshal(ual, usedAddresses); err != nil {
		return nil, err
	}

	return usedAddresses, nil
}

func GetHostAddressList(conf config.Config, token, hostName string) (*HostAddress, error) {
	urlPath := conf.Server.AppID + "/addresses/search_hostname/" + hostName + "/"
	hal, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	hostAddress := new(HostAddress)
	if err := json.Unmarshal(hal, hostAddress); err != nil {
		return nil, err
	}

	return hostAddress, nil
}

func CreateIPAddress(conf config.Config, token, subnetID, hostName, address string, i net.Interface) ([]byte, error) {
	urlPath := conf.Server.AppID + "/addresses/"
	caBody := &createIPAddress{
		SubnetId:    subnetID,
		IP:          address,
		Hostname:    hostName,
		Description: i.Name,
		Mac:         i.HardwareAddr.String(),
	}
	res, err := client.PHPIPAMRequest(conf, "POST", token, "", urlPath, caBody)
	if err != nil {
		return nil, err
	}

	return res, nil
}
