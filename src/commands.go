package main

import (
	"fmt"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/daemon"
	"github.com/0x822a5b87/tiny-docker/src/entity"
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
		args := context.Args()
		containerIds := make([]entity.ContainerId, 0)
		for _, arg := range args {
			containerIds = append(containerIds, entity.ContainerId(arg))
		}
		return daemon.SendStopRequest(conf.StopCommand{ContainerIds: containerIds})
	},
}

var logsCommand = cli.Command{
	Name:  constant.Logs.String(),
	Usage: `Print logs of target container`,
	Flags: []cli.Flag{},
	Action: func(context *cli.Context) error {
		args := context.Args()
		containerIds := make([]entity.ContainerId, 0)
		for _, arg := range args {
			containerIds = append(containerIds, entity.ContainerId(arg))
		}
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

var networkCommand = cli.Command{
	Name:  constant.Network.String(),
	Usage: "Operate networks: create/connect/disconnect/rm",
	Subcommands: cli.Commands{
		newNetworkCreateCommand(),
		newNetworkConnectCommand(),
		newNetworkRmCommand(),
		newNetworkInspectCommand(),
	},
	Action: func(c *cli.Context) error {
		return cli.ShowSubcommandHelp(c)
	},
}

func newNetworkCreateCommand() cli.Command {
	return cli.Command{
		Name:  constant.NetworkCreate.String(),
		Usage: "Create a network",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "subnet",
				Usage:    "Specify subnet for the network (e.g. 172.18.0.0/16)",
				Required: false,
			},
			&cli.StringFlag{
				Name:  "driver",
				Usage: "Specify network driver (only bridge supported)",
				Value: "bridge",
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Specify network name",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			networkName := c.String("name")
			if err := daemon.SendNetworkCreate(networkName); err != nil {
				return fmt.Errorf("failed to create network: %w", err)
			}
			fmt.Printf("Network %s created successfully\n", networkName)
			return nil
		},
	}
}

func newNetworkConnectCommand() cli.Command {
	return cli.Command{
		Name:  constant.NetworkConnect.String(),
		Usage: "Connect a container to a network",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "network",
				Usage:    "Network name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "container",
				Usage:    "Container ID/name",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			//networkName := c.String("network")
			//containerID := c.String("container")
			//
			//if err := daemon.ConnectNetwork(networkName, containerID); err != nil {
			//	return fmt.Errorf("failed to connect network: %w", err)
			//}
			//fmt.Printf("Container %s connected to network %s\n", containerID, networkName)
			return nil
		},
	}
}

func newNetworkRmCommand() cli.Command {
	return cli.Command{
		Name:  constant.NetworkRm.String(),
		Usage: "Remove a network",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Network name",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			networkName := c.String("name")
			if err := daemon.SendNetworkRm(networkName); err != nil {
				return fmt.Errorf("failed to remove network: %w", err)
			}
			fmt.Printf("Network %s removed successfully\n", networkName)
			return nil
		},
	}
}

func newNetworkInspectCommand() cli.Command {
	return cli.Command{
		Name:  constant.NetworkInspect.String(),
		Usage: "Inspect a network",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Network name",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			networkName := c.String("name")
			if err := daemon.SendNetworkInspect(networkName); err != nil {
				return fmt.Errorf("failed to remove network: %w", err)
			}
			return nil
		},
	}
}
