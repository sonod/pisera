package model

import (
	"encoding/json"
	"strings"

	"github.com/sonod/pisera/config"
	"github.com/sonod/pisera/lib"
)

type Subnet struct {
	Data []subnetData
}

type subnetData struct {
	Id             string
	Subnet         string
	Mask           string
	SectionId      string
	Description    string
	MasterSubnetId string
	Usage          struct {
		Used              string
		Maxhosts          string
		Freehosts         string
		Freehosts_percent float64
		Used_percent      float64
	}
}

type UsageSubnet struct {
	Data struct {
		Used             string
		Maxhosts         string
		Freehosts        string
		FreehostsPercent float64
		UsedPercent      float64
	}
}

type createSubnet struct {
	Subnet    string `json:"subnet"`
	Mask      string `json:"mask"`
	SectionId string `json:"sectionId"`
}

func GetSubnetList(conf config.Config, token, secid string) (*Subnet, error) {
	urlPath := conf.Server.AppID + "/sections/" + secid + "/subnets/"
	s, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	subnet := new(Subnet)
	if err := json.Unmarshal(s, subnet); err != nil {
		return nil, err
	}

	return subnet, nil
}

func GetSubnetID(conf config.Config, token, secID, cidr string) (string, error) {
	subnet, err := GetSubnetList(conf, token, secID)
	if err != nil {
		return "", nil
	}
	cidrSplit := strings.Split(cidr, "/")
	var snID string
	for _, sn := range subnet.Data {
		if sn.Subnet == cidrSplit[0] && sn.Mask == cidrSplit[1] {
			snID = sn.Id
		}
	}

	return snID, nil
}

func GetUsageSubnet(conf config.Config, token, subid string) (*UsageSubnet, error) {
	urlPath := conf.Server.AppID + "/subnets/" + subid + "/usage/"
	us, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	usageSubnet := new(UsageSubnet)
	if err := json.Unmarshal(us, usageSubnet); err != nil {
		return nil, err
	}

	return usageSubnet, nil
}

func CreateSubnet(conf config.Config, token, cidr, secid string) ([]byte, error) {
	urlPath := conf.Server.AppID + "/subnets/"
	ci := strings.Split(cidr, "/")
	csbody := &createSubnet{
		Subnet:    ci[0],
		Mask:      ci[1],
		SectionId: secid,
	}
	res, err := client.PHPIPAMRequest(conf, "POST", token, "", urlPath, csbody)
	if err != nil {
		return nil, err
	}
	return res, nil
}
