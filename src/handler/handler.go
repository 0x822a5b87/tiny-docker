package handler

import (
	"fmt"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

type ActionHandler func(req Request) (rsp Response, err error)

func AddHandler(action constant.Action, ac ActionHandler) {
	registry[action] = ac
}

func handleRequest(req Request) (Response, error) {
	h, ok := registry[req.Act]
	if !ok {
		return ErrorResponse(fmt.Errorf("%s", req.Act.String()), constant.ErrUnsupportedAction)
	}
	return h(req)
}
