package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/handler"
	"github.com/sirupsen/logrus"
)

func handlePs(request handler.Request) (handler.Response, error) {
	// TODO implement real query.
	return handler.Response{
		Code: 0,
		Msg:  "success",
		Data: []map[string]interface{}{
			{"id": "abc123", "pid": 12345, "status": "running"},
			{"id": "def456", "pid": 67890, "status": "exited"},
		},
	}, nil
}

func handleCommit(request handler.Request) (handler.Response, error) {
	commands, err := handler.ParamsFromRequest[conf.CommitCommands](&request)
	if err != nil {
		logrus.Errorf("error parse request: %s\n", err.Error())
		return handler.ErrorMessageResponse("type convert error", constant.ErrUnsupportedAction)
	}
	err = Commit(commands)
	if err != nil {
		return handler.ErrorResponse(err, constant.ErrMalformedUdsRsp)
	}
	return handler.SuccessResponse(""), nil

}
