package model

import (
	"encoding/json"

	"github.com/sonod/pisera/config"
	"github.com/sonod/pisera/lib"
)

type Section struct {
	Data []sectionData
}

type sectionData struct {
	Id               string
	Name             string
	Description      string
	MasterSection    string
	Permission       string
	StrictMode       string
	SubnetOrdering   string
	Order            string
	EditData         string
	ShowVLAN         string
	ShowVRF          string
	ShowSupernetOnly string
	DNS              string
}

func getSectionList(conf config.Config, token string) ([]byte, error) {
	urlPath := conf.Server.AppID + "/sections/"
	section, err := client.PHPIPAMRequest(conf, "GET", token, "", urlPath, nil)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func GetSectionID(conf config.Config, token, sectionName string) (string, error) {
	sel, err := getSectionList(conf, token)
	if err != nil {
		return "", err
	}

	section := new(Section)
	if err := json.Unmarshal(sel, section); err != nil {
		return "", err
	}

	var secID string
	for _, s := range section.Data {
		if s.Name == sectionName {
			secID = s.Id
		}
	}

	return secID, nil
}
