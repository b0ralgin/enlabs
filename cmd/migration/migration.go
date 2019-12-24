package migration

import (
	"enlabs/config"
	"enlabs/db"
	_ "enlabs/db/migrations"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/urfave/cli"
)

// StartMigrationCommand - старт HTTP сервера для обработки запросов по REST API
func StartMigrationCommand() cli.Command {
	return cli.Command{
		Name:   "migrate",
		Usage:  "start migration",
		Action: startMigration,
	}
}

func startMigration(c *cli.Context) error {
	var cfg config.AppConfig
	if loadCfgErr := config.LoadConfig(c.App.Name, &cfg); loadCfgErr != nil {
		return errors.Wrap(loadCfgErr, "can't load config")
	}
	db, dbErr := db.NewPostgresClient(cfg.DSN)
	if dbErr != nil {
		return errors.Wrap(dbErr, "can't connect to DB")
	}

	err := goose.Up(db.GetConn(), "db/migrations")
	if err != nil {
		return errors.Wrap(err, "can't perfrom migration")
	}
	return nil
}
