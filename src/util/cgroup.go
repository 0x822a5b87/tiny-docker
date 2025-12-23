package util

import (
	"os"
	"path"
)

func GetCgroupPath(name string, cgroupPath string, autoCreate bool) (string, error) {
	var err error
	if _, err = os.Stat(cgroupPath); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err = os.Mkdir(cgroupPath, 0755); err == nil {
			} else {
				return "", err
			}
		}
		realPath := path.Join(cgroupPath, name)
		return realPath, nil
	}
	return "", err
}
