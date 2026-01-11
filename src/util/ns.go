package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/sirupsen/logrus"
)

func NsenterExec(pid int, args []string, env []string) error {
	path, err := exec.LookPath(constant.Nsenter)
	if err != nil {
		panic(err)
	}
	argv := []string{
		constant.Nsenter,
		"-t", fmt.Sprintf("%d", pid),
		"-a",
	}
	argv = append(argv, args...)
	if err = syscall.Exec(path, argv, env); err != nil {
		logrus.Errorf("error exec %s : %v", path, err)
		return err
	}
	return nil
}

func ReadNsEnv(pid int) ([]string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\x00"), nil
}
