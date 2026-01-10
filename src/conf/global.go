package conf

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// GlobalConfig NOTE THAT IT CAN ONLY BE USED IN THE CLIENT CONTEXT. ONCE `Start()` IS CALLED,
// THE SPAWNED PROCESS WILL NOT BE ABLE TO ACCESS THE CONFIG FILE BECAUSE THE FILE SYSTEM HAS BEEN MODIFIED.
var GlobalConfig Config

func LoadDaemonConfig() {
	loadConfig(Commands{})
}

func LoadRunConfig(commands RunCommands) {
	loadConfig(commands.IntoCommands())
}

func LoadCommitConfig(cmd CommitCommands) {
	loadConfig(cmd.IntoCommands())
}

func LoadBasicCommand() {
	loadConfig(Commands{})
}

func loadConfig(commands Commands) {
	loadFile()
	GlobalConfig.Cmd = commands
	environ()
}

func loadFile() {
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

func environ() {
	env := make([]string, 0)
	env = appendEnv(env, MetaName, GlobalConfig.Meta.Name)
	env = appendEnv(env, FsBasePath, GlobalConfig.Fs.Root)

	env = appendEnv(env, FsReadLayerPath, GlobalConfig.ReadPath())
	env = appendEnv(env, FsWriteLayerPath, GlobalConfig.WritePath())
	env = appendEnv(env, FsWorkLayerPath, GlobalConfig.WorkPath())
	env = appendEnv(env, FsMergeLayerPath, GlobalConfig.MergePath())

	env = appendEnv(env, RuntimeDockerdUdsFile, GlobalConfig.DockerdUdsFile())
	env = appendEnv(env, RuntimeDockerdUdsPidFile, GlobalConfig.DockerdUdsPidFile())
	env = appendEnv(env, RuntimeDockerdLogFile, GlobalConfig.DockerdLogFile())
	env = appendEnv(env, RuntimeDockerdContainerStatus, GlobalConfig.DockerdContainerStatusPath())

	if GlobalConfig.Cmd.Detach {
		env = appendEnv(env, DetachMode, "true")
	} else {
		env = appendEnv(env, DetachMode, "false")
	}

	GlobalConfig.InnerEnv = env
}

func appendEnv(environ []string, key EnvVariable, value string) []string {
	return append(environ, fmt.Sprintf("%s=%s", key.String(), value))
}
