package handler

import (
	"encoding/json"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

type Request struct {
	Act    constant.Action `json:"act"`
	Params []byte          `json:"params"`
}

func ParamsIntoRequest[T any](act constant.Action, data T) (*Request, error) {
	params, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &Request{
		Act:    act,
		Params: params,
	}, nil
}

func ParamsFromRequest[T any](request *Request) (T, error) {
	var t T
	err := json.Unmarshal(request.Params, &t)
	return t, err
}
