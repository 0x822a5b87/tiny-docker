package handler

import (
	"fmt"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

type Response struct {
	Code int    `json:"code"` // 0=成功，非0=失败
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func SuccessResponse(data any) Response {
	return Response{
		Code: constant.UdsStatusOk,
		Msg:  "success",
		Data: data,
	}
}

func ErrorMessageResponse(msg string, wrapErr constant.Err) (Response, error) {
	return ErrorResponse(fmt.Errorf("%s", msg), wrapErr)
}

func ErrorResponse(err error, wrapErr constant.Err) (Response, error) {
	return Response{
		Code: wrapErr.ErrorCode,
		Msg:  wrapErr.Wrap(err).Error(),
	}, err
}
