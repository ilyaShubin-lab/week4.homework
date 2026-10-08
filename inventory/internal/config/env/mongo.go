package env

import (
	"fmt"
	"net"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	Port     string `env:"MONGO_PORT,required"`
	Database string `env:"MONGO_DATABASE,required"`
	AuthDB   string `env:"MONGO_AUTH_DB,required"`
	Username string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
}

type mongoConfig struct {
	raw mongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &mongoConfig{raw: raw}, nil
}

// mongodb://inventory_admin:inventory_secret@localhost:27018/?authSource=admin
func (cfg *mongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s/?authSource=%s",
		cfg.raw.Username,
		cfg.raw.Password,
		net.JoinHostPort(cfg.raw.Host, cfg.raw.Port),
		cfg.raw.AuthDB,
	)
}

func (cfg *mongoConfig) DatabaseName() string {
	return cfg.raw.Database
}
