package daemon

import (
	"fmt"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

func Exec(command *conf.ExecCommand) error {
	conf.LoadBasicCommand()
	p := getContainerStatusFilePath(command.Id)
	state, err := readContainerState(p)
	if err != nil {
		logrus.Errorf("exec read container state error: %v", err)
		return err
	}
	if state.Status != entity.ContainerRunning {
		return fmt.Errorf("exec: container %s is not running", command.Id)
	}

	pid := state.Pid
	env, err := util.ReadNsEnv(pid)
	if err != nil {
		logrus.Errorf("exec read container env error: %v", err)
		return err
	}
	return util.NsenterExec(pid, command.Args, env)
}
