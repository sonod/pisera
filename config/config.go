package config

import "github.com/BurntSushi/toml"

type Config struct {
	Server server `toml:"server"`
}

type server struct {
	Endpoint string
	Username string
	Password string
	AppID    string `toml:"app_id"`
}

func NewConfig(configPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)

	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func defaultConfig(c *Config) {
	c.Server.Endpoint = "http://localhost/api/"
}
