package util

import (
	"fmt"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

func GenPidPath(pid int) string {
	return fmt.Sprintf("%s/%s", constant.CgroupBasePath, constant.DefaultContainerName)
}
