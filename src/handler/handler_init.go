package handler

import "github.com/0x822a5b87/tiny-docker/src/constant"

var registry map[constant.Action]ActionHandler

func init() {
	registry = make(map[constant.Action]ActionHandler)
}
