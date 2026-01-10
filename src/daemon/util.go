package daemon

import (
	"github.com/0x822a5b87/tiny-docker/src/conf"
	"github.com/0x822a5b87/tiny-docker/src/util"
	"github.com/sirupsen/logrus"
)

func formatEventOutput() {

}

// setupUnionFsFromEnv init read-layer, write-layer, work-layer, merge-layer for daemon
func setupUnionFsFromEnv() error {
	readPath := conf.FsReadLayerPath.Get()
	writePath := conf.FsWriteLayerPath.Get()
	workPath := conf.FsWorkLayerPath.Get()
	mergePath := conf.FsMergeLayerPath.Get()
	logrus.Infof("read path: {%s}, write path: {%s}, work path: {%s}, merge path : {%s}", readPath, writePath, workPath, mergePath)
	// mount -t overlay overlay -o lowerdir=...,upperdir=...,workdir=... /root/tiny-docker/busybox/merged
	if err := util.MountOverlayFS(readPath, writePath, workPath, mergePath); err != nil {
		logrus.Errorf("mount proc error : %s", err.Error())
		return err
	}
	return nil
}

//func setupUnionFsFromConfig() error {
//	readPath := conf.GlobalConfig.ReadPath()
//	writePath := conf.GlobalConfig.WritePath()
//	workPath := conf.GlobalConfig.WorkPath()
//	mergePath := conf.GlobalConfig.MergePath()
//	logrus.Infof("read path: {%s}, write path: {%s}, work path: {%s}, merge path : {%s}", readPath, writePath, workPath, mergePath)
//	// mount -t overlay overlay -o lowerdir=...,upperdir=...,workdir=... /root/tiny-docker/busybox/merged
//	if err := util.MountOverlayFS(readPath, writePath, workPath, mergePath); err != nil {
//		logrus.Errorf("mount proc error : %s", err.Error())
//		return err
//	}
//	return nil
//}

func clearUnionFsFromConfig() error {
	mergePath := conf.GlobalConfig.MergePath()
	logrus.Infof("clear merge path: %s", mergePath)
	if err := util.UnmountOverlayFS(mergePath); err != nil {
		logrus.Errorf("mount proc error : %s", err.Error())
		return err
	}
	return nil
}
