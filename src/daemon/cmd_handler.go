package daemon

import (
	"encoding/json"

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

	data, err := json.Marshal(containers)
	if err != nil {
		logrus.Errorf("error marshal containers: %s", err.Error())
		return handler.ErrorMessageResponse(err.Error(), constant.ErrExecCommand)
	}

	return handler.SuccessResponse(string(data)), nil
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
	return handler.SuccessResponse(""), nil
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
	return handler.SuccessResponse(""), nil
}

func handleContainerStop(request handler.Request) (handler.Response, error) {
	c, err := handler.ParamsFromRequest[entity.Container](&request)
	if err != nil {
		logrus.Errorf("error parse container stop request: %s", err.Error())
		return handler.ErrorMessageResponse("error parse container stop request", constant.ErrMalformedUdsReq)
	}
	err = stopContainer(c)
	if err != nil {
		return handler.ErrorResponse(err, constant.ErrMalformedUdsRsp)
	}
	return handler.SuccessResponse(""), nil
}
