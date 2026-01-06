package conf

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var GlobalConfig Config

func LoadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logrus.Errorf("Error reading config file: %v", err)
		panic(err)
	}

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		logrus.Errorf("Error unmarshal config file: %v", err)
		panic(err)
	}
}
