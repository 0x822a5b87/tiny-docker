package manager

import (
	"os"

	"github.com/0x822a5b87/tiny-docker/src/subsystem"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
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

	data := []byte(value.Into())
	if err = os.WriteFile(cgroupPath, data, 0644); err != nil {
		logrus.Errorf("error write value : {%s}", string(data))
		return err
	}

	return nil
}
