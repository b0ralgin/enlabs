package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type AppConfig struct {
	DSN      string `envconfig:"dsn"`
	LogLevel string `envconfig:"log_level"`
	Addr     string `envconfig:"addr"`
}

// LoadConfig ...
func LoadConfig(app string, c interface{}) error {
	if proccessErr := envconfig.Process(app, c); proccessErr != nil {
		return errors.Wrapf(proccessErr, "failed to load config for %s", app)
	}
	return nil
}