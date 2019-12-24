package main

import (
	"enlabs/cmd/server"
	"github.com/urfave/cli"
	"os"
)

var Name string

func main() {
	app := cli.NewApp()
	app.Name = "enlabs"
	app.Commands = []cli.Command{
		server.StartHTTPServerCommand(),
	}

	if runErr := app.Run(os.Args); runErr != nil {
		panic(runErr)
	}
}
