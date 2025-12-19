package main

import (
	"errors"
)

// StartContainer macOS 下空实现（仅报错提示）
func StartContainer(cmd string) error {
	return errors.New("StartContainer is not supported on macOS, run on Linux instead")
}
