package server

import (
	"enlabs/config"
	"enlabs/db"
	"enlabs/pkg/account"
	"enlabs/pkg/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// StartHTTPServerCommand - старт HTTP сервера для обработки запросов по REST API
func StartHTTPServerCommand() cli.Command {
	return cli.Command{
		Name:   "server",
		Usage:  "Start HTTP REST service",
		Action: startRouter,
	}
}

func startRouter(c *cli.Context) error {
	var cfg config.AppConfig
	if loadCfgErr := config.LoadConfig(c.App.Name, &cfg); loadCfgErr != nil {
		return errors.Wrap(loadCfgErr, "can't load confifg")
	}
	logLevel, parseErr := logrus.ParseLevel(cfg.LogLevel)
	if parseErr != nil {
		return errors.Wrap(parseErr, "can't parse log level")
	}
	log := logrus.New()
	log.SetLevel(logLevel)

	db, dbErr := db.NewPostgresClient(cfg.DSN)
	if dbErr != nil {
		return errors.Wrap(dbErr, "can't connect to DB")
	}
	am := account.NewAccountManager(db)
	hs := api.NewHTTPServer(am, log.WithField("api", "http"))
	if err := hs.Run(cfg.Addr); err != nil {
		return errors.Wrap(err, "http server has stopped")
	}
	return nil
}
