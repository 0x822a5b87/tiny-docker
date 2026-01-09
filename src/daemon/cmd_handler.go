package daemon

//func handleRun(request handler.Request) (handler.Response, error) {
//	commands, err := handler.ParamsFromRequest[conf.RunCommands](&request)
//	if err != nil {
//		logrus.Errorf("error parse request: %s\n", err.Error())
//		return handler.ErrorMessageResponse("type convert error", constant.ErrUnsupportedAction)
//	}
//	err = RunContainer(commands)
//	if err != nil {
//		return handler.ErrorResponse(err, constant.ErrMalformedUdsRsp)
//	}
//	return handler.SuccessResponse(""), nil
//}
