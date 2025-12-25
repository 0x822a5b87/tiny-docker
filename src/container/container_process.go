package container

import (
	"os"
	"os/exec"

	"github.com/0x822a5b87/tiny-docker/src/constant"
	"golang.org/x/sys/unix"
)

func NewParentProcess(tty bool, commands []string) *exec.Cmd {
	args := []string{"init"}
	for _, command := range commands {
		args = append(args, command)
	}
	cmd := exec.Command(constant.UnixProcSelfExe, args...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUTS |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWNS |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWIPC,
		Unshareflags: unix.CLONE_NEWNS,
		Setpgid:      true,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
