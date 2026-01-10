package conf

import (
	"encoding/json"
	"path/filepath"

	"github.com/0x822a5b87/tiny-docker/src/constant"
)

type PathType string

var ImagePath PathType = "images"
var RuntimePath PathType = "runtime"
var LogPath PathType = "logs"
var StatePath PathType = "state"
var ContainerPath PathType = "container"

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

type PsCommand struct {
	All bool
}

type StopCommand struct {
	ContainerIds []string
}

type LogsCommand struct {
	ContainerId string
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

func (c Config) String() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) ImageName() string {
	return ExtractNameFromTarPath(c.Cmd.Image)
}

func (c Config) ReadPath() string {
	return filepath.Join(c.rootPath(), string(ImagePath), "read", c.ImageName())
}

func (c Config) WritePath() string {
	return c.buildIndPath(ImagePath, "write", c.Cmd.Id)
}

func (c Config) WorkPath() string {
	return c.buildIndPath(ImagePath, "work", c.Cmd.Id)
}

func (c Config) MergePath() string {
	return c.buildIndPath(ImagePath, "merge", c.Cmd.Id)
}

func (c Config) DockerdUdsFile() string {
	return c.DockerdPath(RuntimePath, constant.DockerdUdsConnFile)
}

func (c Config) DockerdUdsPidFile() string {
	return c.DockerdPath(RuntimePath, constant.DockerdUdsPidFile)
}

func (c Config) DockerdLogFile() string {
	return c.DockerdPath(LogPath, constant.DockerdLogFile)
}

func (c Config) DockerdContainerStatusPath() string {
	return c.DockerdPath(StatePath, "")
}

func (c Config) DockerdContainerLogPath() string {
	return c.DockerdPath(ContainerPath, "")
}

func (c Config) rootPath() string {
	root := c.Fs.Root
	if c.Cmd.Volume != "" {
		root = c.Cmd.Volume
	}
	return filepath.Join(root)
}

func (c Config) DockerdPath(pathType PathType, fileName string) string {
	if fileName == "" {
		return filepath.Join(c.rootPath(), string(pathType))
	}
	return filepath.Join(c.rootPath(), string(pathType), fileName)
}

func (c Config) ImagePath(pathType PathType) string {
	return filepath.Join(c.rootPath(), string(pathType), c.ImageName())
}

func (c Config) buildSharePath(pathType PathType, suffix string) string {
	return filepath.Join(c.rootPath(), string(pathType), suffix)
}

func (c Config) buildIndPath(pathType PathType, suffix string, id string) string {
	return filepath.Join(c.rootPath(), string(pathType), suffix, id)
}
