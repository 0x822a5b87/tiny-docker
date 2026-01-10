package daemon

import (
	"strings"
	"time"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/handler"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

// NOTE THAT ALL CLIENT EVENT CAN ONLY BE INVOKED IN CLIENT

func SendPsRequest(command conf.PsCommand) error {
	conf.LoadBasicCommand()
	return sendRequest(constant.Ps, command)
}

func SendCommitRequest(commands conf.CommitCommands) error {
	return sendRequest[conf.CommitCommands](constant.Commit, commands)
}

func SendContainerInitRequest(pid int) error {
	c := entity.Container{
		Id:        conf.GlobalConfig.Cmd.Id,
		Pid:       pid,
		Image:     conf.GlobalConfig.ImageName(),
		Command:   strings.Join(conf.GlobalConfig.Cmd.Args, " "),
		CreatedAt: time.Now().UnixMilli(),
		Status:    entity.ContainerRunning,
		Name:      conf.GlobalConfig.ImageName(),
	}

	return sendRequest(constant.Run, c)
}

func SendContainerExitRequest() error {
	c := entity.Container{
		Id:     conf.GlobalConfig.Cmd.Id,
		ExitAt: time.Now().UnixMilli(),
	}
	return sendRequest(constant.Stop, c)
}

func sendRequest[D any](act constant.Action, data D) error {
	req, err := handler.ParamsIntoRequest[D](act, data)
	if err != nil {
		return err
	}
	err, rsp := handler.SendRequest(req)
	if err != nil {
		logrus.Errorf("error sending commit request: %v\n", err)
		return err
	}
	util.LogWithoutExtraInfo(rsp)
	return nil
}
