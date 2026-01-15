package util

import (
	"os"
	"strings"

	"github.com/0x822a5b87/tiny-docker/src/conf"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetBoolEnv(key string) bool {
	v := os.Getenv(key)
	return strings.ToLower(v) == "true"
}

func InitTestConfig() {
	conf.GlobalConfig = conf.Config{
		Meta: conf.MetaConfig{Name: "test"},
		Fs:   conf.FsConfig{Root: "/root/test-tiny-docker/"},
	}
	conf.Environ()
}
