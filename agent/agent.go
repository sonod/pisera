package agent

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/sonod/pisera/config"
	"github.com/sonod/pisera/model"
)

func Run(conf config.Config, section string) error {
	for {
		token, err := model.GetAuthenticationToken(conf)
		if err != nil {
			return err
		}

		sectionID, err := model.GetSectionID(conf, token, section)
		if err != nil {
			return err
		}

		h, err := os.Hostname()
		if err != nil {
			panic(err)
		}

		sl, err := model.GetSubnetList(conf, token, sectionID)
		if err != nil {
			return err
		}

		subnetsList := map[string]bool{}
		for _, s := range sl.Data {
			subnetsList[s.Subnet+"/"+s.Mask] = true
		}

		interfaces, err := net.Interfaces()
		if err != nil {
			panic(err)
		}

		for _, i := range interfaces {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
					if ip.IP.To4() != nil {
						_, cidr, _ := net.ParseCIDR(addr.String())
						if _, v := subnetsList[cidr.String()]; v == false {
							_, err := model.CreateSubnet(conf, token, cidr.String(), sectionID)
							if err != nil {
								panic(err)
							}
						}
						subnetID, err := model.GetSubnetID(conf, token, sectionID, cidr.String())
						if err != nil {
							panic(err)
						}

						ual, err := model.GetUsedAddressList(conf, token, subnetID)
						if err != nil {
							panic(err)
						}

						userAddressList := map[string]bool{}
						for _, ua := range ual.Data {
							userAddressList[ua.IP] = true
						}

						a := strings.Split(addr.String(), "/")
						if _, v := userAddressList[a[0]]; v == false {
							if _, err := model.CreateIPAddress(conf, token, subnetID, h, a[0], i); err != nil {
								panic(err)
							}
						}
					}
				}
			}
		}
		time.Sleep(time.Duration(1) * time.Hour)
	}
	return nil
}