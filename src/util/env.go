package util

import (
	"os"
	"strings"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetBoolEnv(key string) bool {
	v := os.Getenv(key)
	return strings.ToLower(v) == "true"
}
