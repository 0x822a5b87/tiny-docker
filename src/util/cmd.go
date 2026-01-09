package util

import (
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func IsDaemonProcess(args []string) bool {
	return len(args) == 1 && args[0] == constant.Daemon.String()
}

func GetInitArgs(context *cli.Context) ([]string, error) {
	commands := make([]string, 0)
	if len(context.Args()) < 1 {
		logrus.Error("GetInitArgs error : No command to execute")
		return commands, constant.ErrMalformedArgs
	}

	for _, arg := range context.Args() {
		commands = append(commands, arg)
	}
	return commands, nil
}

func GetImageAndArgs(context *cli.Context) (string, []string, error) {
	commands := make([]string, 0)
	if len(context.Args()) < 1 {
		logrus.Error("GetImageAndArgs error : No command to execute")
		return "", commands, constant.ErrMalformedArgs
	}

	image := context.Args().First()
	for i := 1; i < context.NArg(); i++ {
		commands = append(commands, context.Args().Get(i))
	}

	return image, commands, nil
}
