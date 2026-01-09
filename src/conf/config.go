package conf

import (
	"path/filepath"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

type Commands struct {
	Id       string
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
	fullID, _ := GenContainerID()
	return Commands{
		Id:      fullID,
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
	Meta     MetaConfig `yaml:"meta"`
	Fs       FsConfig   `yaml:"fs"`
	Cmd      Commands   `yaml:"cmd"`
	InnerEnv []string   `yaml:"inner_env"`
}

type MetaConfig struct {
	Name string `yaml:"name"`
}

type FsConfig struct {
	Root string `yaml:"root"`
}

func (c Config) ImageName() string {
	return ExtractNameFromTarPath(c.Cmd.Image)
}

func (c Config) ReadPath() string {
	return c.buildSharePath("images", "read")
}

func (c Config) WritePath() string {
	return c.buildIndPath("images", "write", c.Cmd.Id)
}

func (c Config) WorkPath() string {
	return c.buildIndPath("images", "work", c.Cmd.Id)
}

func (c Config) MergePath() string {
	return c.buildIndPath("images", "merge", c.Cmd.Id)
}

func (c Config) DockerdUdsFile() string {
	return filepath.Join(c.RootPath("runtime"), constant.DockerdUdsConnFile)
}

func (c Config) DockerdUdsPidFile() string {
	return filepath.Join(c.RootPath("runtime"), constant.DockerdUdsPidFile)
}

func (c Config) DockerdLogFile() string {
	return filepath.Join(c.RootPath("logs"), constant.DockerdLogFile)
}

func (c Config) RootPath(pathType string) string {
	root := c.Fs.Root
	if c.Cmd.Volume != "" {
		root = c.Cmd.Volume
	}
	return filepath.Join(root, pathType, c.ImageName())
}

func (c Config) buildSharePath(pathType string, suffix string) string {
	return filepath.Join(c.RootPath(pathType), suffix)
}

func (c Config) buildIndPath(pathType, suffix string, id string) string {
	return filepath.Join(c.RootPath(pathType), suffix, id)
}
