package conf

import (
	"os"
	"strings"
)

type EnvVariable string

func (v EnvVariable) String() string {
	return string(v)
}

func (v EnvVariable) Get() string {
	return v.getValue()
}

func (v EnvVariable) GetBoolean() bool {
	value := v.getValue()
	return strings.ToLower(value) == "true"
}

func (v EnvVariable) getValue() string {
	value := os.Getenv(v.String())
	if value != "" {
		return value
	}
	if len(GlobalConfig.InnerEnv) > 0 {
		for _, k := range GlobalConfig.InnerEnv {
			parts := strings.SplitN(k, "=", 2)
			if len(parts) == 2 && parts[0] == v.String() {
				return parts[1]
			}
		}
	}
	return ""
}

const MetaName EnvVariable = "tiny-docker-meta-name"

const FsBasePath EnvVariable = "tiny-docker-fs-base-path"

const FsReadLayerPath EnvVariable = "tiny-docker-fs-read-layer-path"
const FsWriteLayerPath EnvVariable = "tiny-docker-fs-write-layer-path"
const FsWorkLayerPath EnvVariable = "tiny-docker-fs-work-layer-path"
const FsMergeLayerPath EnvVariable = "tiny-docker-fs-merge-layer-path"

const DetachMode EnvVariable = "tiny-docker-detach-mode"

const RuntimeDockerdUdsFile EnvVariable = "tiny-docker-runtime-dockerd-uds-file"
const RuntimeDockerdUdsPidFile EnvVariable = "tiny-docker-runtime-dockerd-pid-file"
const RuntimeDockerdLogFile EnvVariable = "tiny-docker-runtime-dockerd-log-file"
