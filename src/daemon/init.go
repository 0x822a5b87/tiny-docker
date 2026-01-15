package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/handler"
	"github.com/0x822a5b87/tiny-docker/src/network"
	"github.com/sirupsen/logrus"
)

var networks *network.Networks

func init() {
	handler.AddHandler(constant.Ps, handlePs)
	handler.AddHandler(constant.Commit, handleCommit)
	handler.AddHandler(constant.Run, handleContainerRun)
	handler.AddHandler(constant.Stop, handleContainerStop)
	handler.AddHandler(constant.Logs, handleContainerLogs)
	handler.AddHandler(constant.Wait, handleWaitContainer)

	handler.AddHandler(constant.NetworkCreate, handleNetworkCreate)
	handler.AddHandler(constant.NetworkRm, handleNetworkRm)
	handler.AddHandler(constant.NetworkInspect, handleNetworkInspect)
	handler.AddHandler(constant.NetworkConnect, handleNetworkCreate)

	conf.LoadBasicCommand()
	var err error
	networks, err = network.NewNetworks()
	if err != nil {
		logrus.Errorf("error creating networks: %v", err)
		panic(err)
	}
}
