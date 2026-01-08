package conf

import (
	"path/filepath"

	"github.com/0x822a5b87/tiny-docker/src/util"
)

type Commands struct {
	Tty      bool
	Detach   bool
	Image    string
	DstImage string
	Args     []string
	Cfg      CgroupConfig
	UserEnv  []string
	Volume   string
}

type RunCommands struct {
	Tty     bool
	Detach  bool
	Image   string
	Args    []string
	Cfg     CgroupConfig
	UserEnv []string
	Volume  string
}

func (r RunCommands) IntoCommands() Commands {
	return Commands{
		Tty:     r.Tty,
		Detach:  r.Detach,
		Image:   r.Image,
		Args:    r.Args,
		Cfg:     r.Cfg,
		UserEnv: r.UserEnv,
		Volume:  r.Volume,
	}
}

type CommitCommands struct {
	SrcName string
	DstName string
	Volume  string
}

func (c CommitCommands) IntoCommands() Commands {
	return Commands{
		Image:    c.SrcName,
		DstImage: c.DstName,
		Volume:   c.Volume,
	}
}

type CgroupConfig struct {
	MemoryLimit string
	CpuShares   string
}

type Config struct {
	Meta MetaConfig `yaml:"meta"`
	Fs   FsConfig   `yaml:"fs"`
	Cmd  Commands
}

type MetaConfig struct {
	Name string `yaml:"name"`
}

type FsConfig struct {
	Root string `yaml:"root"`
}

func (c Config) ImageName() string {
	return util.ExtractNameFromTarPath(c.Cmd.Image)
}

func (c Config) ReadPath() string {
	return c.buildPath("read")
}

func (c Config) WritePath() string {
	return c.buildPath("write")
}

func (c Config) WorkPath() string {
	return c.buildPath("work")
}

func (c Config) MergePath() string {
	return c.buildPath("merge")
}

func (c Config) RootPath() string {
	root := c.Fs.Root
	if c.Cmd.Volume != "" {
		root = c.Cmd.Volume
	}
	return filepath.Join(root, c.ImageName())
}

func (c Config) buildPath(suffix string) string {
	return filepath.Join(c.RootPath(), suffix)
}
