package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/handler"
)

func init() {
	addAllHandler()
}

func addAllHandler() {
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

}
