package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sonod/pisera/agent"
	"github.com/sonod/pisera/config"
	"github.com/sonod/pisera/model"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version     bool
		agentEnable bool
		configPath  string
		section     string
		cidr        string
		hostname    string
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n  %s [OPTIONS] ARGS...\nArgs\n", os.Args[0], os.Args[0])
		fmt.Fprint(os.Stderr, "  subnet-list  Subnet List\n")
		fmt.Fprint(os.Stderr, "  free-address First Address List(require -cidr)\n")
		fmt.Fprint(os.Stderr, "  address-list Used Address List(require -cidr)\n")
		fmt.Fprint(os.Stderr, "  usage-subnet Usage Subnet(require -cidr)\nOptions\n")
		flags.PrintDefaults()
	}
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&agentEnable, "agent", false, "Pisera_Agent_Mode")
	flags.StringVar(&section, "section", "Customers", "PHPIPAM_Section")
	flags.StringVar(&cidr, "cidr", "", "your cidr(ex: 172.16.0.0/24)")
	flags.StringVar(&configPath, "config", "/etc/pisera.toml", "Pisera_Config")
	flags.StringVar(&hostname, "hostname", "", "Server Address List")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	conf, err := config.NewConfig(configPath)
	if err != nil {
		return ExitCodeError
	}

	if agentEnable {
		if err := agent.Run(conf, section); err != nil {
			return ExitCodeError
		}
		return ExitCodeOK
	}

	token, err := model.GetAuthenticationToken(conf)
	if err != nil {
		panic(err)
	}

	if len(hostname) != 0 {
		hal, err := model.GetHostAddressList(conf, token, hostname)
		if err != nil {
			panic(err)
		}
		hostAddressTable := tablewriter.NewWriter(os.Stdout)
		hostnameList := []string{"hostname"}
		interfaceList := []string{hostname}
		for _, ha := range hal.Data {
			hostnameList = append(hostnameList, ha.Description)
			interfaceList = append(interfaceList, ha.IP)
		}
		hostAddressTable.SetHeader(hostnameList)
		hostAddressTable.Append(interfaceList)
		hostAddressTable.Render()
		return ExitCodeOK
	}

	sectionID, err := model.GetSectionID(conf, token, section)
	if err != nil {
		panic(err)
	}

	if flags.Arg(0) == "subnet-list" {
		sl, err := model.GetSubnetList(conf, token, sectionID)
		if err != nil {
			return ExitCodeError
		}
		subnetListTable := tablewriter.NewWriter(os.Stdout)
		subnetListTable.SetHeader([]string{"Subnet ID", "Subnet List"})
		for _, s := range sl.Data {
			if s.Subnet != "0.0.0.0" {
				subnetListTable.Append([]string{s.Id, s.Subnet + "/" + s.Mask})
			}
		}
		subnetListTable.Render()
		return ExitCodeOK
	} else {
		if cidr == "" {
			return ExitCodeError
		}
		subnetID, err := model.GetSubnetID(conf, token, sectionID, cidr)
		if err != nil {
			return ExitCodeError
		}
		switch flags.Arg(0) {
		case "free-address":
			fa, err := model.GetFirstFreeAddress(conf, token, subnetID)
			if err != nil {
				panic(err)
			}
			freeAddressTable := tablewriter.NewWriter(os.Stdout)
			freeAddressTable.SetHeader([]string{"Free IPAddress"})
			freeAddressTable.Append([]string{fa.Data})
			freeAddressTable.Render()
			return ExitCodeOK
		case "address-list":
			ual, err := model.GetUsedAddressList(conf, token, subnetID)
			if err != nil {
				panic(err)
			}
			usedAddressTable := tablewriter.NewWriter(os.Stdout)
			usedAddressTable.SetHeader([]string{"Hostname", "IP Address", "Mac Address", "Description"})
			for _, ua := range ual.Data {
				usedAddressTable.Append([]string{ua.Hostname, ua.IP, ua.Mac, ua.Description})
			}
			usedAddressTable.Render()
			return ExitCodeOK
		case "usage-subnet":
			us, err := model.GetUsageSubnet(conf, token, subnetID)
			if err != nil {
				panic(err)
			}
			usageSubnetTable := tablewriter.NewWriter(os.Stdout)
			usageSubnetTable.SetHeader([]string{"Max Address", "Used Addresses", "Free Addresses"})
			usageSubnetTable.Append([]string{us.Data.Maxhosts, us.Data.Used, us.Data.Freehosts})
			usageSubnetTable.Render()
			return ExitCodeOK
		default:
			return ExitCodeError
		}
	}

	return ExitCodeOK
}
