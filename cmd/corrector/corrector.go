package corrector

import (
	"enlabs/config"
	"enlabs/db"
	"enlabs/pkg/account"
	"enlabs/pkg/scheduler"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// StartCorrectorCommand - старт HTTP сервера для обработки запросов по REST API
func StartCorrectorCommand() cli.Command {
	return cli.Command{
		Name:   "corrector",
		Usage:  "Start Corrector service",
		Action: startCorrector,
	}
}

func startCorrector(c *cli.Context) error {
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
	ct := account.NewCorrector(db)
	return scheduler.RunScheduler(ct, cfg.Period)
}
