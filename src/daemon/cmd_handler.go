package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/entity"
	"github.com/0x822a5b87/tiny-docker/src/handler"
	"github.com/sirupsen/logrus"
)

func handlePs(request handler.Request) (handler.Response, error) {
	command, err := handler.ParamsFromRequest[conf.PsCommand](&request)
	if err != nil {
		logrus.Errorf("error parse request: %s", err.Error())
		return handler.ErrorMessageResponse("type convert error", constant.ErrMalformedUdsReq)
	}

	containers, err := ps(command)
	if err != nil {
		logrus.Errorf("error parse request: %s", err.Error())
		return handler.ErrorMessageResponse(err.Error(), constant.ErrExecCommand)
	}

	return handler.SuccessResponse(containers)
}

func handleCommit(request handler.Request) (handler.Response, error) {
	commands, err := handler.ParamsFromRequest[conf.CommitCommands](&request)
	if err != nil {
		logrus.Errorf("error parse request: %s", err.Error())
		return handler.ErrorMessageResponse("type convert error", constant.ErrMalformedUdsReq)
	}
	err = Commit(commands)
	if err != nil {
		return handler.ErrorResponse(err, constant.ErrMalformedUdsRsp)
	}
	return handler.SuccessResponse("{}")
}

func handleContainerRun(request handler.Request) (handler.Response, error) {
	c, err := handler.ParamsFromRequest[entity.Container](&request)
	if err != nil {
		logrus.Errorf("error parse container run request: %s", err.Error())
		return handler.ErrorMessageResponse("error parse status", constant.ErrMalformedUdsReq)
	}
	err = runContainer(c)
	if err != nil {
		return handler.ErrorResponse(err, constant.ErrMalformedUdsRsp)
	}
	return handler.SuccessResponse("{}")
}

func handleContainerStop(request handler.Request) (handler.Response, error) {
	containers, err := handler.ParamsFromRequest[[]entity.Container](&request)
	if err != nil {
		logrus.Errorf("error parse container stop request: %s", err.Error())
		return handler.ErrorMessageResponse("error parse container stop request", constant.ErrMalformedUdsReq)
	}
	err = stopContainers(containers)
	if err != nil {
		return handler.ErrorResponse(err, constant.ErrMalformedUdsRsp)
	}
	return handler.SuccessResponse("{}")
}

func handleContainerLogs(request handler.Request) (handler.Response, error) {
	container, err := handler.ParamsFromRequest[entity.Container](&request)
	if err != nil {
		logrus.Errorf("error parse container stop request: %s", err.Error())
		return handler.ErrorMessageResponse("error parse container stop request", constant.ErrMalformedUdsReq)
	}

	data, err := logs(container.Id)
	if err != nil {
		return handler.ErrorResponse(err, constant.ErrMalformedUdsRsp)
	}
	return handler.SuccessResponse(data)
}

func handleWaitContainer(request handler.Request) (handler.Response, error) {
	waitReq, err := handler.ParamsFromRequest[entity.WaitRequest](&request)
	if err != nil {
		logrus.Errorf("error parse wait container request: %s", err.Error())
		return handler.ErrorMessageResponse("error parse container stop request", constant.ErrMalformedUdsReq)
	}
	go wait(waitReq)
	return handler.SuccessResponse("{}")
}
