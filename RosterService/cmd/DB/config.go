package DB

import "github.com/kelseyhightower/envconfig"

const appID = "UnitService"

type Config struct {
	DBType           string `envconfig:"database_type" default:"mysql"`
	DBName           string `envconfig:"database_name" default:"roster_db"`
	DBUsername       string `envconfig:"database_username" default:"root"`
	DBPassword       string `envconfig:"database_password" default:"Future1994!)"`
	DBMigrationsPath string `envconfig:"database_migrations_path" default:"file://db/migrations"`
}

func ParseEnvDB() *Config {
	c := new(Config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil
	}
	return c
}
