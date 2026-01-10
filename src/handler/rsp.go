package handler

import (
	"encoding/json"
	"fmt"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

type Response struct {
	Code int    `json:"code"` // 0=成功，非0=失败
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func DataIntoResponse[T any](code int, msg string, data T) (Response, error) {
	rspDataBytes, err := json.Marshal(data)
	if err != nil {
		return Response{}, err
	}
	rspDataStr := string(rspDataBytes)
	return Response{
		Code: code,
		Msg:  msg,
		Data: rspDataStr,
	}, nil
}

func DataFromResponse[T any](rsp Response) (T, error) {
	var t T
	var dataBytes []byte

	switch v := rsp.Data.(type) {
	case string:
		var rawJSON string
		err := json.Unmarshal([]byte(v), &rawJSON)
		if err != nil {
			dataBytes = []byte(v)
		} else {
			dataBytes = []byte(rawJSON)
		}
	case []byte:
		dataBytes = v
	default:
		return t, fmt.Errorf("unsupported data type: %T, expected string or []byte", v)
	}

	err := json.Unmarshal(dataBytes, &t)
	if err != nil {
		s, ok := rsp.Data.(T)
		if !ok {
			return t, fmt.Errorf("failed to unmarshal to %T: %v, raw data: %s", t, err, string(dataBytes))
		}
		return s, nil
	}
	return t, nil
}

func SuccessResponse(data any) (Response, error) {
	return DataIntoResponse(constant.UdsStatusOk, "success", data)
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
