package main

import (
	"enlabs/cmd/corrector"
	"enlabs/cmd/migration"
	"enlabs/cmd/server"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "enlabs"
	app.Commands = []cli.Command{
		server.StartHTTPServerCommand(),
		migration.StartMigrationCommand(),
		corrector.StartCorrectorCommand(),
	}

	if runErr := app.Run(os.Args); runErr != nil {
		panic(runErr)
	}
}
