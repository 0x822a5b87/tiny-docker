package main

import (
	"os"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/container"
	"github.com/0x822a5b87/tiny-docker/src/util"
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
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "c",
			Usage: "cpu share limit",
		},
		cli.StringFlag{
			Name:  "entrypoint",
			Usage: "Overwrite the default ENTRYPOINT of the image",
		},
		&cli.StringSliceFlag{
			Name:  "env,e",
			Usage: "Set environment variables, format: `KEY=VALUE`",
		},
		&cli.StringSliceFlag{
			Name:  "volume,v",
			Usage: "Set volumes for the container",
		},
	},
	Action: func(context *cli.Context) error {
		image, args, err := util.GetImageAndArgs(context)
		if err != nil {
			log.Error(err, "error parse image and args")
			return err
		}
		runCommands := conf.RunCommands{}
		runCommands.Tty = context.Bool("it")
		runCommands.Detach = context.Bool("d")
		runCommands.Volume = context.String("v")
		runCommands.Image = image
		runCommands.Args = args
		runCommands.Cfg = conf.CgroupConfig{
			MemoryLimit: context.String("m"),
			CpuShares:   context.String("c"),
		}
		runCommands.UserEnv = context.StringSlice("env")
		return container.Run(runCommands)
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: `Init container process run user's process in container. Do not call it outside.`,
	Action: func(context *cli.Context) error {
		log.Infof("init come on pid : %d", os.Getpid())
		args, err := util.GetArgs(context)
		if err != nil {
			return err
		}
		return container.RunContainerInitProcess(args[0], args)
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: `Create a compression file(.tar) from a container`,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "s",
			Usage: "Name of container to commit",
		},
		&cli.StringSliceFlag{
			Name:  "t",
			Usage: "Target name of committed container",
		},
		&cli.StringSliceFlag{
			Name:  "v",
			Usage: "Volume of container to commit",
		},
	},
	Action: func(context *cli.Context) error {
		cmd := conf.CommitCommands{
			SrcName: context.String("s"),
			DstName: context.String("t"),
			Volume:  context.String("v"),
		}
		return container.Commit(cmd)
	},
}
