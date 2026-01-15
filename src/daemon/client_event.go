package daemon

import (
	"fmt"
	"strings"
	"time"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/handler"
	"github.com/sirupsen/logrus"
)

// NOTE THAT ALL CLIENT EVENT CAN ONLY BE INVOKED IN CLIENT

func SendPsRequest(command conf.PsCommand) error {
	conf.LoadBasicCommand()
	rsp, err := sendRequest[conf.PsCommand](constant.Ps, command)
	if err != nil {
		return err
	}
	if rsp.Code != constant.UdsStatusOk {
		return fmt.Errorf(rsp.Msg)
	}
	data, err := handler.DataFromResponse[[]entity.Container](*rsp)
	if err != nil {
		return err
	}
	formatContainerTable(data)
	return nil
}

func SendStopRequest(command conf.StopCommand) error {
	conf.LoadBasicCommand()
	containers := make([]entity.Container, 0)
	for _, id := range command.ContainerIds {
		containers = append(containers, entity.Container{
			Id:     id,
			ExitAt: time.Now().UnixMilli(),
		})
	}
	rsp, err := sendRequest[[]entity.Container](constant.Stop, containers)
	if err != nil {
		return err
	}
	if rsp.Code != constant.UdsStatusOk {
		return fmt.Errorf(rsp.Msg)
	}
	return nil
}

func SendStopCurrentRequest() error {
	containers := make([]entity.Container, 0)
	c := entity.Container{
		Id:     conf.GlobalConfig.Cmd.Id,
		ExitAt: time.Now().UnixMilli(),
	}
	containers = append(containers, c)
	_, err := sendRequest(constant.Stop, containers)
	return err
}

func SendCommitRequest(commands conf.CommitCommands) error {
	conf.LoadCommitConfig(commands)
	_, err := sendRequest[conf.CommitCommands](constant.Commit, commands)
	if err != nil {
		return err
	}
	return nil
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

	_, err := sendRequest(constant.Run, c)
	return err
}

func SendLogRequest(command conf.LogsCommand) error {
	conf.LoadBasicCommand()
	rsp, err := sendRequest[entity.Container](constant.Logs, entity.Container{Id: command.ContainerId})
	if err != nil {
		return err
	}
	if rsp.Code != constant.UdsStatusOk {
		return fmt.Errorf(rsp.Msg)
	}

	logStr, err := handler.DataFromResponse[string](*rsp)
	if err != nil {
		return err
	}
	logStr = strings.ReplaceAll(logStr, "\\n", "\n")
	logStr = strings.ReplaceAll(logStr, `\"`, `"`)
	logStr = strings.Trim(logStr, `"`)
	logStr = strings.TrimSpace(logStr)
	fmt.Println(logStr)
	return nil
}

func SendNetworkCreate(name string) error {
	conf.LoadBasicCommand()
	rsp, err := sendRequest[entity.Network](constant.NetworkCreate, entity.Network{Name: name})
	if err != nil {
		return err
	}
	if rsp.Code != constant.UdsStatusOk {
		return fmt.Errorf(rsp.Msg)
	}

	network, err := handler.DataFromResponse[string](*rsp)
	if err != nil {
		return err
	}
	fmt.Println(network)
	return nil
}

func SendNetworkRm(name string) error {
	conf.LoadBasicCommand()
	rsp, err := sendRequest[entity.Network](constant.NetworkRm, entity.Network{Name: name})
	if err != nil {
		return err
	}
	if rsp.Code != constant.UdsStatusOk {
		return fmt.Errorf(rsp.Msg)
	}

	network, err := handler.DataFromResponse[string](*rsp)
	if err != nil {
		return err
	}
	fmt.Println(network)
	return nil
}

func SendNetworkInspect(name string) error {
	conf.LoadBasicCommand()
	rsp, err := sendRequest[entity.Network](constant.NetworkInspect, entity.Network{Name: name})
	if err != nil {
		return err
	}
	if rsp.Code != constant.UdsStatusOk {
		return fmt.Errorf(rsp.Msg)
	}

	network, err := handler.DataFromResponse[string](*rsp)
	if err != nil {
		return err
	}
	fmt.Println(network)
	return nil
}

func SendWaitRequest(request entity.WaitRequest) error {
	_, err := sendRequest(constant.Wait, request)
	return err
}

func sendRequest[D any](act constant.Action, data D) (*handler.Response, error) {
	req, err := handler.ParamsIntoRequest[D](act, data)
	if err != nil {
		return nil, err
	}
	err, rsp := handler.SendRequest(req)
	if err != nil {
		logrus.Errorf("error sending commit request: %v\n", err)
		return nil, err
	}
	return &rsp, nil
}
