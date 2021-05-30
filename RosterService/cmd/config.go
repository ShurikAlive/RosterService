package cmd

import (
	"github.com/kelseyhightower/envconfig"
)

const appID = "UnitService"

type Config struct {
	ServeRESTAddress string `envconfig:"servr_rest_address" default:":8080"`
	UnitServiceBasePath string `http://localhost:8181/`
	EquipmentServiceBasePath string `http://localhost:8181/`
}

func ParseEnv() (*Config, error) {
	c := new(Config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, err
	}
	return c, nil
}
