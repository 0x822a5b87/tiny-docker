package manager

import (
	"os"

	"github.com/0x822a5b87/tiny-docker/src/util"
)

type CgroupFileSystem struct {
	Path       string
	AutoCreate bool
}

func NewCgroupFileSystem(basePath string, autoCreate bool) *CgroupFileSystem {
	return &CgroupFileSystem{Path: basePath, AutoCreate: autoCreate}
}

func (f *CgroupFileSystem) Read(name string) (error, string) {
	cgroupPath, err := util.GetCgroupPath(name, f.Path, f.AutoCreate)
	if err != nil {
		return err, ""
	}

	data, err := os.ReadFile(cgroupPath)
	if err != nil {
		return err, ""
	}

	return nil, string(data)
}

func (f *CgroupFileSystem) Write(name string, data string) error {
	cgroupPath, err := util.GetCgroupPath(name, f.Path, f.AutoCreate)
	if err != nil {
		return err
	}
	return os.WriteFile(cgroupPath, []byte(data), 0644)
}
