package util

import (
	"fmt"
	"os"
	"os/exec"
)

func AppendEnv(cmd *exec.Cmd, key, value string) {
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
