package ns

import (
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

// StartContainer Linux 下启动带 UTS Namespace 的容器
func StartContainer(cmd string) error {
	command := exec.Command(cmd)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Linux 专属：UTS Namespace 配置
	command.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS,
	}

	return command.Run()
}
