package main

import (
	"fmt"
	"os"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/container"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit tiny-docker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "c",
			Usage: "cpu share limit",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("it")
		cfg := conf.CgroupConfig{
			MemoryLimit: context.String("m"),
			CpuShares:   context.String("c"),
		}
		return Run(tty, cmd, cfg)
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `Init container process run user's process in container. Do not call it outside.`,
	Action: func(context *cli.Context) error {
		log.Infof("init come on pid : %d", os.Getpid())
		cmd := context.Args().Get(0)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
