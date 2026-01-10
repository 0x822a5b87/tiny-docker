package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/handler"
)

func init() {
	handler.AddHandler(constant.Ps, handlePs)
	handler.AddHandler(constant.Commit, handleCommit)
	handler.AddHandler(constant.Run, handleContainerRun)
	handler.AddHandler(constant.Stop, handleContainerStop)
	handler.AddHandler(constant.Logs, handleContainerLogs)
}
