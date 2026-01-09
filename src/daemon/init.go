package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/handler"
)

func init() {
	handler.AddHandler(constant.Ps, handlePs)
}
