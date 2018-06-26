# pisera(phpIPAM Server Assistant)

## Description
pisera is the CLI tool and the Agent tool of phpIPAM.
In the CLI, you can work with phpIPAM to return a list of addresses or return unused addresses.
Agent automatically registers the host's address in phpIPAM.
Also, if the server has IPMI, it registers the device based on the IPMI address.

## Usage
setting phpIPAM endpoint and username and password and app_id to `pisera.toml`.
default path is `/etc/pisera.toml`.

### Client Mode
check subnet-list
`$ pisera subnet-list`

check device-list
`$ pisera device-list`

check free-address for subnet
`$ pisera -subnet 172.17.1.0/24 free-address`

check address-list for subnet
`$ pisera -subnet 172.17.1.0/24 address-list`

check usage-subnet
`$ pisera -subnet 172.17.1.0/24 usage-subnet`

check address for hostname
`$ pisera -hostname nrm`

### Agent Mode
launch pisera agent
`$ pisera -agent`

## Install

To install, use `go get github.com/sonod/pisera`:

## ToDo

add test

## Contribution

1. Fork ([https://github.com/sonod/pisera/fork](https://github.com/sonod/pisera/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[sonod](https://github.com/sonod)
