package conf

import (
	"path/filepath"

	"github.com/0x822a5b87/tiny-docker/src/util"
)

type RunCommands struct {
	Tty     bool
	Image   string
	Args    []string
	Cfg     CgroupConfig
	UserEnv []string
}

func (r RunCommands) ImageName() string {
	return util.ExtractNameFromTarPath(r.Image)
}

type CgroupConfig struct {
	MemoryLimit string
	CpuShares   string
}

type Config struct {
	Meta   MetaConfig `yaml:"meta"`
	Fs     FsConfig   `yaml:"fs"`
	RunCmd RunCommands
}

type MetaConfig struct {
	Name string `yaml:"name"`
}

type FsConfig struct {
	Root string `yaml:"root"`
}

func (c Config) ReadPath() string {
	return filepath.Join(c.Fs.Root, c.RunCmd.ImageName(), "read")
}

func (c Config) WritePath() string {
	return filepath.Join(c.Fs.Root, c.RunCmd.ImageName(), "write")
}

func (c Config) WorkPath() string {
	return filepath.Join(c.Fs.Root, c.RunCmd.ImageName(), "work")
}

func (c Config) MergePath() string {
	return filepath.Join(c.Fs.Root, c.RunCmd.ImageName(), "merge")
}
