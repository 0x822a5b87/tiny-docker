package daemon

import (
	"os/exec"

	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/constant"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

func setupEnv(cmd *exec.Cmd) {
	util.AppendEnv(cmd, constant.MetaName, conf.GlobalConfig.Meta.Name)
	util.AppendEnv(cmd, constant.FsBasePath, conf.GlobalConfig.Fs.Root)

	util.AppendEnv(cmd, constant.FsReadLayerPath, conf.GlobalConfig.ReadPath())
	util.AppendEnv(cmd, constant.FsWriteLayerPath, conf.GlobalConfig.WritePath())
	util.AppendEnv(cmd, constant.FsWorkLayerPath, conf.GlobalConfig.WorkPath())
	util.AppendEnv(cmd, constant.FsMergeLayerPath, conf.GlobalConfig.MergePath())

	if conf.GlobalConfig.Cmd.Detach {
		util.AppendEnv(cmd, constant.DetachMode, "true")
	} else {
		util.AppendEnv(cmd, constant.DetachMode, "false")
	}
}

// setupUnionFsFromEnv init read-layer, write-layer, work-layer, merge-layer for daemon
func setupUnionFsFromEnv() error {
	readPath := util.GetEnv(constant.FsReadLayerPath)
	writePath := util.GetEnv(constant.FsWriteLayerPath)
	workPath := util.GetEnv(constant.FsWorkLayerPath)
	mergePath := util.GetEnv(constant.FsMergeLayerPath)
	logrus.Infof("read path: {%s}, write path: {%s}, work path: {%s}, merge path : {%s}", readPath, writePath, workPath, mergePath)
	// mount -t overlay overlay -o lowerdir=...,upperdir=...,workdir=... /root/tiny-docker/busybox/merged
	if err := util.MountOverlayFS(readPath, writePath, workPath, mergePath); err != nil {
		logrus.Errorf("mount proc error : %s", err.Error())
		return err
	}
	return nil
}

func setupUnionFsFromConfig() error {
	readPath := conf.GlobalConfig.ReadPath()
	writePath := conf.GlobalConfig.WritePath()
	workPath := conf.GlobalConfig.WorkPath()
	mergePath := conf.GlobalConfig.MergePath()
	logrus.Infof("read path: {%s}, write path: {%s}, work path: {%s}, merge path : {%s}", readPath, writePath, workPath, mergePath)
	// mount -t overlay overlay -o lowerdir=...,upperdir=...,workdir=... /root/tiny-docker/busybox/merged
	if err := util.MountOverlayFS(readPath, writePath, workPath, mergePath); err != nil {
		logrus.Errorf("mount proc error : %s", err.Error())
		return err
	}
	return nil
}

func clearUnionFsFromConfig() error {
	mergePath := conf.GlobalConfig.MergePath()
	logrus.Infof("clear merge path: %s", mergePath)
	if err := util.UnmountOverlayFS(mergePath); err != nil {
		logrus.Errorf("mount proc error : %s", err.Error())
		return err
	}
	return nil
}
