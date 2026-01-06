package conf

import (
	"path/filepath"
)

type CgroupConfig struct {
	MemoryLimit string
	CpuShares   string
}

type Config struct {
	Meta MetaConfig `yaml:"meta"`
	Fs   FsConfig   `yaml:"fs"`
}

type MetaConfig struct {
	Name string `yaml:"name"`
}

type FsConfig struct {
	Root string `yaml:"root"`
}

func (c Config) ReadPath() string {
	return filepath.Join(c.Fs.Root, c.Meta.Name, "read")
}

func (c Config) WritePath() string {
	return filepath.Join(c.Fs.Root, c.Meta.Name, "write")
}

func (c Config) WorkPath() string {
	return filepath.Join(c.Fs.Root, c.Meta.Name, "work")
}

func (c Config) MergePath() string {
	return filepath.Join(c.Fs.Root, c.Meta.Name, "merge")
}
