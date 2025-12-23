package manager

import (
	"os"

	"github.com/0x822a5b87/tiny-docker/src/subsystem"
	"github.com/0x822a5b87/tiny-docker/src/util"
)

func Write[I subsystem.Item, V subsystem.Value](f *CgroupFileSystem, ss subsystem.Subsystem[I, V]) error {
	cgroupPath, err := util.GetCgroupPath(ss.Name(), f.Path, f.AutoCreate)
	if err != nil {
		return err
	}

	value, err := ss.Get()
	if err != nil {
		return err
	}

	if err := os.WriteFile(cgroupPath, []byte(value.Into()), 0644); err != nil {
		return err
	}

	return nil
}
