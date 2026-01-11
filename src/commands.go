package main

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/daemon"
	"github.com/0x822a5b87/tiny-docker/src/util"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var daemonCommand = cli.Command{
	Name:  constant.Daemon.String(),
	Usage: `Start dockerd to serve all API requests`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug mode",
		},
	},
	Action: func(context *cli.Context) error {
		debug := context.Bool("debug")
		return daemon.StartDockerd(debug)
	},
}

var initCommand = cli.Command{
	Name:  constant.InitDaemon.String(),
	Usage: `Init daemon process run user's process in daemon. Do not call it outside.`,
	Action: func(context *cli.Context) error {
		return daemon.RunDaemon()
	},
}

var runCommand = cli.Command{
	Name:  string(constant.Run),
	Usage: `Create a daemon with namespace and cgroups limit tiny-docker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach daemon",
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
			Usage: "Set volumes for the daemon",
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
		err = daemon.RunContainerCmd(runCommands)
		if err != nil {
			log.Errorf("error sending run request: %v\n", err)
			return err
		}
		return nil
	},
}

var initContainerCommand = cli.Command{
	Name:  constant.InitContainer.String(),
	Usage: `Init container process run user's process in daemon. Do not call it outside.`,
	Action: func(context *cli.Context) error {
		args, err := util.GetInitArgs(context)
		if err != nil {
			return err
		}
		return daemon.RunContainer(args[0], args)
	},
}

var commitCommand = cli.Command{
	Name:  constant.Commit.String(),
	Usage: `Create a compression file(.tar) from a daemon`,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "s",
			Usage: "Id of daemon to commit",
		},
		&cli.StringSliceFlag{
			Name:  "t",
			Usage: "Target name of committed daemon",
		},
		&cli.StringSliceFlag{
			Name:  "v",
			Usage: "Volume of daemon to commit",
		},
	},
	Action: func(context *cli.Context) error {
		cmd := conf.CommitCommands{
			SrcName: context.String("s"),
			DstName: context.String("t"),
			Volume:  context.String("v"),
		}

		return daemon.SendCommitRequest(cmd)
	},
}

var psCommand = cli.Command{
	Name:  constant.Ps.String(),
	Usage: `List target containers`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "a",
			Usage: "Show all containers",
		},
	},
	Action: func(context *cli.Context) error {
		all := context.Bool("a")
		return daemon.SendPsRequest(conf.PsCommand{All: all})
	},
}

var stopCommand = cli.Command{
	Name:  constant.Stop.String(),
	Usage: `Stop target container`,
	Flags: []cli.Flag{},
	Action: func(context *cli.Context) error {
		containerIds := context.Args()
		return daemon.SendStopRequest(conf.StopCommand{ContainerIds: containerIds})
	},
}

var logsCommand = cli.Command{
	Name:  constant.Logs.String(),
	Usage: `Print logs of target container`,
	Flags: []cli.Flag{},
	Action: func(context *cli.Context) error {
		containerIds := context.Args()
		if len(containerIds) != 1 {
			return constant.ErrMalformedLogsArgs
		}
		return daemon.SendLogRequest(conf.LogsCommand{ContainerId: containerIds[0]})
	},
}

var execCommand = cli.Command{
	Name:  constant.Exec.String(),
	Usage: `Execute a command in a running container`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	Action: func(context *cli.Context) error {
		command, err := util.GetExecArgs(context)
		if err != nil {
			return err
		}
		return daemon.Exec(command)
	},
}
