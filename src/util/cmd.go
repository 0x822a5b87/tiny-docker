package util

import (
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/urfave/cli"
)

func GetArgs(context *cli.Context) ([]string, error) {
	commands := make([]string, 0)
	if len(context.Args()) < 1 {
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
		return "", commands, constant.ErrMalformedArgs
	}

	image := context.Args().First()
	for i := 1; i < context.NArg(); i++ {
		commands = append(commands, context.Args().Get(i))
	}

	return image, commands, nil
}
