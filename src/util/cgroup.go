package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

func GenContainerCgroupPath(id string) string {
	cgroupBasePath := constant.CgroupBasePath
	cgroupServicePath := filepath.Join(cgroupBasePath, constant.CgroupServiceName)
	if err := EnsureFilePathExist(cgroupServicePath); err != nil {
		panic(err)
	}
	controllers := []string{"cpu", "memory", "io"}
	err := EnableCgroupControllers(controllers, cgroupServicePath)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s", cgroupServicePath, id)
}

// EnableCgroupControllers edit cgroup.subtree_control so that we can use cpu.max to manage CPU resources.
func EnableCgroupControllers(controllers []string, cgroupBasePath string) error {
	subtreeFile := filepath.Join(cgroupBasePath, constant.CgroupSubtreeControl)
	current, err := os.ReadFile(subtreeFile)
	if err != nil {
		return fmt.Errorf("read subtree control failed: %w", err)
	}
	currentStr := strings.TrimSpace(string(current))
	currentControllers := strings.Fields(currentStr)

	var enableCmd []string
	for _, c := range controllers {
		exists := false
		for _, item := range currentControllers {
			if strings.TrimPrefix(item, "+") == c && strings.TrimPrefix(item, "-") == c {
				exists = true
				break
			}
		}
		if !exists {
			enableCmd = append(enableCmd, "+"+c)
		}
	}

	if len(enableCmd) == 0 {
		return nil
	}

	cmdStr := strings.Join(enableCmd, " ") + "\n"
	f, err := os.OpenFile(subtreeFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open subtree control file failed: %w", err)
	}
	defer func() { _ = f.Close() }()

	_, err = f.WriteString(cmdStr)
	if err != nil {
		return fmt.Errorf("write subtree control failed: %w", err)
	}

	return nil
}

func GetCgroupPath(name string, cgroupPath string, autoCreate bool) (string, error) {
	var err error
	if _, err = os.Stat(cgroupPath); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(cgroupPath, 0755); err == nil {
			} else {
				return "", err
			}
		}
		realPath := path.Join(cgroupPath, name)
		return realPath, nil
	}
	return "", err
}
