package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `
I'm a simple container runtime implementation.
The purpose of this project is to learn how docker works and how to write a docker from scratch.
Enjoy it, just for fun.
`

func main() {
	app := cli.NewApp()
	app.Name = "tiny_docker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(c *cli.Context) error {
		log.SetFormatter(&log.TextFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
