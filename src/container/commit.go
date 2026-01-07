package container

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

func Commit(cmd conf.CommitCommands) error {
	logrus.Infof("Commit Commands: %s", cmd)
	conf.LoadCommitConfig(cmd)
	if err := SetupUnionFsFromConfig(); err != nil {
		logrus.Error("error setting up union fs", err)
		return err
	}
	defer func() {
		err := ClearUnionFsFromConfig()
		if err != nil {
			logrus.Error("error clear up union fs", err)
		}
	}()
	root := conf.GlobalConfig.MergePath()
	return util.Tar(conf.GlobalConfig.Cmd.DstImage, root)
}
