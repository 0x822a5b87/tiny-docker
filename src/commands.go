package main

import (
	"os"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/container"
	"github.com/0x822a5b87/tiny-docker/src/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type RunCommands struct {
	Tty      bool
	Commands []string
	Cfg      conf.CgroupConfig
	UserEnv  []string
}

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

		&cli.StringSliceFlag{
			Name:  "env,e",
			Usage: "Set environment variables, format: `KEY=VALUE`",
		},
	},
	Action: func(context *cli.Context) error {
		commands, err := util.GetCommands(context)
		if err != nil {
			return err
		}
		runCommands := RunCommands{}
		runCommands.Tty = context.Bool("it")
		runCommands.Cfg = conf.CgroupConfig{
			MemoryLimit: context.String("m"),
			CpuShares:   context.String("c"),
		}
		runCommands.UserEnv = context.StringSlice("env")
		runCommands.Commands = commands
		return Run(runCommands)
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `Init container process run user's process in container. Do not call it outside.`,
	Action: func(context *cli.Context) error {
		log.Infof("init come on pid : %d", os.Getpid())
		args, err := util.GetCommands(context)
		if err != nil {
			return err
		}
		return container.RunContainerInitProcess(args[0], args)
	},
}
